package manager

import (
	"net/rpc"
	"sync"
)

// Manager ...
type Manager struct {
	sync.RWMutex
	rooms map[*Room]struct{}
}

var manager = &Manager{rooms: map[*Room]struct{}{}}

// Close ...
func Close() {
	manager.Lock()
	defer manager.Unlock()
	for r := range manager.rooms {
		r.Close()
	}
}

func Remove(r *Room) {
	manager.Lock()
	defer manager.Unlock()
	delete(manager.rooms, r)
}

// Enter ...
func Enter(c *rpc.Client) *Room {
	manager.Lock()
	defer manager.Unlock()
	for r := range manager.rooms {
		// TODO: room choice by ranking
		if !r.Enter(c) {
			continue
		}
		return r
	}
	r := NewRoom()
	manager.rooms[r] = struct{}{}
	r.Enter(c)
	return r
}
