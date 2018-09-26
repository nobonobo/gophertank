package schema

// Action ...
type Action interface {
	action()
}

// ActionBase ...
type ActionBase struct{}

func (a *ActionBase) action() {}
