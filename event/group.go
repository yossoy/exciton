package event

// Group is interface for event group
type Group interface {
	AddHandler(name string, handler Handler) error
	AddHandlerWithResult(name string, handler HandlerWithResult) error
	AddGroup(name string) (Group, error)
	//	SetUnmatchedHandler(handler UnmatchedHandler)
}
