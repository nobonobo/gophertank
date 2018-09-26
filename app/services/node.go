package services

import (
	"net/rpc"

	"github.com/nobonobo/gophertank/app/actions"
	"github.com/nobonobo/gophertank/app/dispatcher"
	"github.com/nobonobo/gophertank/app/schema"
	"github.com/nobonobo/gophertank/app/store"
)

func init() {
	rpc.Register(&Node{})
}

// Node ...
type Node struct{}

// GetIdentity ...
func (n *Node) GetIdentity(req *struct{}, res *schema.Identity) error {
	*res = store.Identity
	return nil
}

// UpdateMembers ...
func (n *Node) UpdateMembers(req []schema.Identity, res *struct{}) error {
	dispatcher.Dispatch(&actions.UpdateMembers{Members: req})
	return nil
}

// Begin ...
func (n *Node) Begin(req []schema.Identity, res *struct{}) error {
	dispatcher.Dispatch(&actions.Begin{Members: req})
	return nil
}

// Abort ...
func (n *Node) Abort(req *struct{}, res *struct{}) error {
	dispatcher.Dispatch(&actions.Abort{})
	return nil
}

// End ...
func (n *Node) End(req *struct{}, res *struct{}) error {
	dispatcher.Dispatch(&actions.End{})
	return nil
}

// CreateOffer ...
func (n *Node) CreateOffer(uuid string, sdp *string) error {
	act := &actions.CreateOffer{
		UUID:     uuid,
		Response: make(chan string),
	}
	dispatcher.Dispatch(act)
	*sdp = <-act.Response
	return nil
}

// CreateAnswer ...
func (n *Node) CreateAnswer(offer string, answer *string) error {
	act := &actions.CreateAnswer{
		SDP:      offer,
		Response: make(chan string),
	}
	dispatcher.Dispatch(act)
	*answer = <-act.Response
	return nil
}

// CreateConn ...
func (n *Node) CreateConn(info *schema.PeerInfo, res *struct{}) error {
	act := &actions.CreateConn{
		UUID: info.UUID,
		SDP:  info.SDP,
	}
	dispatcher.Dispatch(act)
	return nil
}
