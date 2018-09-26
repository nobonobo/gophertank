package manager

import (
	"net/rpc"
	"sort"
	"sync"
	"time"

	"github.com/nobonobo/gophertank/app/schema"
	"github.com/nobonobo/gophertank/app/store"
	"github.com/nobonobo/gophertank/logger"
)

type State int

const (
	// RollCall 出欠確認
	RollCall State = iota
	// Preperation 準備中
	Preperation
	// Closed Roomの役割終了
	Closed
)

type member struct {
	*rpc.Client
	id schema.Identity
}

// Room ...
type Room struct {
	sync.RWMutex
	members  map[*member]struct{}
	modified time.Time
	state    State
	begin    sync.Once
}

// NewRoom ...
func NewRoom() *Room {
	r := &Room{
		members: map[*member]struct{}{},
	}
	return r
}

// State ...
func (r *Room) State() State {
	r.RLock()
	defer r.RUnlock()
	return r.state
}

// Close ...
func (r *Room) Close() {
	r.Lock()
	defer r.Unlock()
	r.state = Closed
	for m := range r.members {
		m.Close()
	}
	Remove(r)
}

// Enter ...
func (r *Room) Enter(c *rpc.Client) bool {
	r.Lock()
	defer r.Unlock()
	if r.state != RollCall {
		return false
	}
	defer func() { r.modified = time.Now() }()
	if len(r.members) >= store.MaxMembers {
		return false
	}
	var res schema.Identity
	if err := c.Call("Node.GetIdentity", schema.None, &res); err != nil {
		logger.Error(err)
		if _, ok := err.(rpc.ServerError); !ok {
			// not remote error
			return false
		}
		return false
	}
	r.members[&member{Client: c, id: res}] = struct{}{}
	logger.Print("enter:", res.Name)
	return true
}

// Leave ...
func (r *Room) Leave(c *rpc.Client) {
	r.Lock()
	defer r.Unlock()
	if r.state != RollCall {
		return
	}
	defer func() { r.modified = time.Now() }()
	for m := range r.members {
		if m.Client == c {
			delete(r.members, m)
			m.Close()
			logger.Print("leave:", m.id.Name)
			return
		}
	}
}

// Members ...
func (r *Room) Members() []schema.Identity {
	r.RLock()
	defer r.RUnlock()
	res := []schema.Identity{}
	for m := range r.members {
		res = append(res, m.id)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res
}

// IsStandBy ...
func (r *Room) IsStandBy() bool {
	r.Lock()
	defer r.Unlock()
	if r.modified.IsZero() || len(r.members) < 2 {
		return false
	}
	if len(r.members) == store.MaxMembers {
		if time.Since(r.modified) < 27*time.Second {
			r.modified = time.Now().Add(-27 * time.Second)
		}
	}
	return time.Since(r.modified) > 30*time.Second
}

// Begin ...
func (r *Room) Begin() {
	r.begin.Do(func() {
		members := map[*member]struct{}{}
		r.Lock()
		r.state = Preperation
		for k, v := range r.members {
			members[k] = v
		}
		r.Unlock()
		logger.Println("preperation begin")
		defer logger.Println("preperation end")
		if !preperation(members) {
			for m := range members {
				go func(p *member) {
					if err := p.Call("Node.Abort", schema.None, schema.None); err != nil {
						logger.Print(p.id.Name, p.id.UUID, err)
					}
				}(m)
			}
		}
	})
}

func preperation(members map[*member]struct{}) bool {
	res := []schema.Identity{}
	for m := range members {
		res = append(res, m.id)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	for m := range members {
		go func(p *member) {
			if err := p.Call("Node.Begin", res, schema.None); err != nil {
				logger.Print(p.id.Name, p.id.UUID, err)
			}
		}(m)
	}
	// TODO: implement 1 term proc.
	for self := range members {
		for target := range members {
			if self == target {
				continue
			}
			logger.Print("connect:", self.id, "->", target.id)
			var offer, answer string
			if err := self.Call("Node.CreateOffer", target.id.UUID, &offer); err != nil {
				logger.Print("dc connect failed:", err)
				return false
			}
			if err := target.Call("Node.CreateAnswer", offer, &answer); err != nil {
				logger.Print("dc connect failed:", err)
				return false
			}
			info := schema.PeerInfo{
				UUID: target.id.UUID,
				SDP:  answer,
			}
			if err := self.Call("Node.CreateConn", &info, schema.None); err != nil {
				logger.Print("dc connect failed:", err)
				return false
			}
		}
	}
	time.Sleep(time.Second)
	var wg sync.WaitGroup
	for m := range members {
		wg.Add(1)
		go func(p *member) {
			defer wg.Done()
			if err := p.Call("Node.End", schema.None, schema.None); err != nil {
				logger.Print(p.id.Name, p.id.UUID, err)
			}
		}(m)
	}
	wg.Wait()
	return true
}
