package store

import (
	"net"
	"net/rpc"

	"github.com/gopherjs/gopherjs/js"
	"github.com/nobonobo/gophertank/app/schema"
)

const (
	// MaxMembers ...
	MaxMembers = 8
)

// Client ...
type Client struct {
	*rpc.Client
	PeerConnection *js.Object
}

var (
	conn net.Conn
	// Identity ...
	Identity schema.Identity
	// Wanted ...
	Wanted bool
	// CurrentRoomMembers ...
	CurrentRoomMembers = []schema.Identity{}
	// Others ...
	Others = map[string]*Client{}
)
