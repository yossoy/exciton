package event

// Result is type of event result
type Result interface {
	Value() Value
	Error() error
}

type eventResult struct {
	value Value
	err   error
}

func (er eventResult) Value() Value {
	return er.value
}
func (er eventResult) Error() error {
	return er.err
}

// NewValueResult create Result with Value.
func NewValueResult(value Value) Result {
	return eventResult{value: value}
}

// NewErrorResult create Result with error.
func NewErrorResult(err error) Result {
	return eventResult{err: err}
}
