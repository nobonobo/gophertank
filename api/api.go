package api

import (
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"golang.org/x/net/websocket"

	"github.com/nobonobo/gophertank/api/manager"
	"github.com/nobonobo/gophertank/app/schema"
	"github.com/nobonobo/gophertank/logger"
)

// API ...
type API struct {
	http.Handler
}

// New ...
func New() http.Handler {
	a := &API{}
	a.Handler = websocket.Handler(a.handle)
	return a
}

func (a *API) handle(conn *websocket.Conn) {
	//conn.PayloadType = websocket.BinaryFrame
	logger.Print("connect:", conn.RemoteAddr())
	defer logger.Print("disconnect:", conn.RemoteAddr())
	c := jsonrpc.NewClient(conn)
	defer c.Close()
	room := manager.Enter(c)
	defer room.Leave(c)
	tick := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-tick.C:
			if room.IsStandBy() {
				room.Begin()
				return
			}
			if room.State() == manager.RollCall {
				if err := c.Call("Node.UpdateMembers", room.Members(), schema.None); err != nil {
					logger.Error(err)
					if _, ok := err.(rpc.ServerError); !ok {
						// not remote error
						return
					}
				}
			}
		}
	}
}
