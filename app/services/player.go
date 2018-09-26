package services

import (
	"net/rpc"

	"github.com/nobonobo/gophertank/app/schema"
	"github.com/nobonobo/gophertank/app/store"
)

func init() {
	rpc.Register(&Player{})
}

// Player ...
type Player struct{}

// GetIdentity ...
func (p *Player) GetIdentity(req *struct{}, res *schema.Identity) error {
	*res = store.Identity
	return nil
}
