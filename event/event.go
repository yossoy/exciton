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
type Handler func(e *Event) error

// HandlerWithResult is internal event handler function type with result
type HandlerWithResult func(e *Event, callback ResponceCallback)

type Signaller interface {
	Core() *Signal
	Emit(Value) error
}

type Slotter interface {
	Core() *Slot
	Bind(h Handler)
	BindWithResult(h HandlerWithResult)
	IsEnabled() bool
	SetValidateEnabledHandler(validator func(name string) bool)
}
