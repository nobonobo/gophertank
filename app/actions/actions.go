package actions

import "github.com/nobonobo/gophertank/app/schema"

// Enter to the room
type Enter struct {
	schema.ActionBase
	Name string
}

// Leave from the room
type Leave struct {
	schema.ActionBase
}

// UpdateMembers ...
type UpdateMembers struct {
	schema.ActionBase
	Members []schema.Identity
}

// Begin ...
type Begin struct {
	schema.ActionBase
	Members []schema.Identity
}

// Abort ...
type Abort struct {
	schema.ActionBase
}

// End ...
type End struct {
	schema.ActionBase
}

// CreateOffer ...
type CreateOffer struct {
	schema.ActionBase
	UUID     string // target UUID
	Response chan string
}

// CreateAnswer ...
type CreateAnswer struct {
	schema.ActionBase
	SDP      string
	Response chan string
}

// CreateConn ...
type CreateConn struct {
	schema.ActionBase
	UUID string // target UUID
	SDP  string // answer
}
