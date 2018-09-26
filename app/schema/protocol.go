package schema

import "fmt"

// None ...
var None = &struct{}{}

// Identity ...
type Identity struct {
	Name string
	UUID string
}

func (id Identity) String() string {
	return fmt.Sprintf("%s(%s)", id.Name, id.UUID)
}

// PeerInfo ...
type PeerInfo struct {
	UUID string
	SDP  string
}
