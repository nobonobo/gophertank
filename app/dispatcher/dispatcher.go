package dispatcher

import "github.com/nobonobo/gophertank/app/schema"

// ID is a unique identifier representing a registered callback function.
type ID int

var idCounter ID
var callbacks = make(map[ID]func(action schema.Action))

// Dispatch dispatches the given action to all registered callbacks.
func Dispatch(action schema.Action) {
	for _, c := range callbacks {
		c(action)
	}
}

// Register registers the callback to handle dispatched actions, the returned
// ID may be used to unregister the callback later.
func Register(callback func(action schema.Action)) ID {
	idCounter++
	id := idCounter
	callbacks[id] = callback
	return id
}

// Unregister unregisters the callback previously registered via a call to
// Register.
func Unregister(id ID) {
	delete(callbacks, id)
}
