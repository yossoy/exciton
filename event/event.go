package event

// Event is internal event argument
type Event struct {
	Name     string
	Argument Value
	Result   interface{}
	Target   EventTarget
	Host     EventHost
}

// ResponceCallback is internal response handler
type ResponceCallback func(result Result)

// Handler is internal event handler function type
type Handler func(e *Event)

// HandlerWithResult is internal event handler function type with result
type HandlerWithResult func(e *Event, callback ResponceCallback)
