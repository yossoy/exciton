package event

import (
	"errors"
)

// Event is internal event argument
type Event struct {
	Name     string
	Argument Value
	Result   interface{}
	Params   map[string]string
}

// ResponceCallback is internal response handler
type ResponceCallback func(result Result)

// Handler is internal event handler function type
type Handler func(e *Event)

// HandlerWithResult is internal event handler function type with result
type HandlerWithResult func(e *Event, callback ResponceCallback)

// UnmatchedHandler is called at unmatched path
type UnmatchedHandler func(path string, parameter map[string]string, callback ResponceCallback)

type command int

const (
	add command = iota
	remove
	group
	emit
	emitAsync
	quit
)

type cmd struct {
	command          command
	name             string
	respCallback     ResponceCallback
	unmatchedHandler UnmatchedHandler
	argument         interface{}
	targetRouter     *Router
	resultChan       chan interface{}
}

func makeResultChan() chan interface{} {
	//TODO: use sync.Pool?
	return make(chan interface{})
}

type eventMgr struct {
	router    *Router
	eventChan chan cmd
}

// Group is interface for event group
type Group interface {
	AddHandler(name string, handler Handler) error
	AddHandlerWithResult(name string, handler HandlerWithResult) error
	AddGroup(name string) (Group, error)
	SetUnmatchedHandler(handler UnmatchedHandler)
}
type eventGroup struct {
	subRouter *Router
}

type rootGroup struct {
}

func (rg *rootGroup) AddHandler(name string, handler Handler) error {
	return AddHandler(name, handler)
}
func (rg *rootGroup) AddHandlerWithResult(name string, handler HandlerWithResult) error {
	return AddHandlerWithResult(name, handler)
}
func (rg *rootGroup) AddGroup(name string) (Group, error) {
	return AddGroup(name)
}
func (rg *rootGroup) SetUnmatchedHandler(handler UnmatchedHandler) {
	SetUnmatchedHandler(handler)
}

func RootGroup() Group {
	return &rootGroup{}
}

func (em *eventMgr) addHandler(tr *Router, name string, handler interface{}) error {
	rc := makeResultChan()

	em.eventChan <- cmd{
		command:      add,
		name:         name,
		argument:     handler,
		resultChan:   rc,
		targetRouter: tr,
	}
	r := <-rc
	if r == nil {
		return nil
	}
	return r.(error)
}

func (em *eventMgr) addGroup(tr *Router, name string) (Group, error) {
	rc := makeResultChan()
	if tr == nil {
		tr = em.router
	}
	em.eventChan <- cmd{
		command:      group,
		name:         name,
		targetRouter: tr,
		resultChan:   rc,
	}
	r := <-rc
	if r == nil {
		panic("invalid responce")
	}
	if e, ok := r.(error); ok {
		return nil, e
	}
	return r.(*eventGroup), nil
}

func (em *eventMgr) emit(name string, argument Value, respCallback ResponceCallback) error {
	rc := makeResultChan()
	em.eventChan <- cmd{
		command:      emit,
		name:         name,
		argument:     argument,
		resultChan:   rc,
		respCallback: respCallback,
	}
	r := <-rc
	if r == nil {
		return nil
	}
	return r.(error)
}

func (em *eventMgr) quit() {
	rc := makeResultChan()
	em.eventChan <- cmd{
		command:    quit,
		resultChan: rc,
	}
	<-rc
	close(em.eventChan)
}

func (em *eventMgr) eventMain() {
	for {
		select {
		case c, ok := <-em.eventChan:
			if !ok {
				break
			}
			switch c.command {
			case add:
				err := c.targetRouter.Add(c.name, c.argument)
				c.resultChan <- err
			case remove:
				//TODO:
			case emit:
				r, params, evn, err := em.router.Match(c.name)
				if err != nil {
					c.resultChan <- err
				} else if r == nil {
					c.resultChan <- errors.New("event not found")
				} else {
					var arg Value
					if c.argument != nil {
						arg = c.argument.(Value)
					}
					e := &Event{
						Name:     evn,
						Argument: arg,
						Params:   params,
					}
					if u, ok := r.(UnmatchedRouteItem); ok {
						uh := u.Item().(UnmatchedHandler)
						go uh(u.PathSegments(), params, c.respCallback)
					} else if h, ok := r.Item().(Handler); ok {
						if c.respCallback != nil {
							c.resultChan <- errors.New("call EmitWithCallback by non-result event")
						} else {
							go h(e)
							c.resultChan <- nil
						}
					} else if h, ok := r.Item().(HandlerWithResult); ok {
						go func() {
							h(e, c.respCallback)
						}()
						c.resultChan <- nil
					} else {
						panic(errors.New("invalid handler type"))
					}
				}
			case group:
				gr, err := c.targetRouter.AddRoute(c.name, NewRouter())
				if err != nil {
					c.resultChan <- err
				} else {
					c.resultChan <- &eventGroup{subRouter: gr}
				}
			case quit:
				c.resultChan <- nil
				return
			}
		}
	}
}

func (eg *eventGroup) AddHandler(name string, handler Handler) error {
	return masterEventMgr.addHandler(eg.subRouter, name, handler)
}

func (eg *eventGroup) AddHandlerWithResult(name string, handler HandlerWithResult) error {
	return masterEventMgr.addHandler(eg.subRouter, name, handler)
}

func (eg *eventGroup) AddGroup(name string) (Group, error) {
	return masterEventMgr.addGroup(eg.subRouter, name)
}

func (eg *eventGroup) SetUnmatchedHandler(handler UnmatchedHandler) {
	eg.subRouter.SetUnmatchedItem(handler)
}

var (
	masterEventMgr *eventMgr
)

// AddHandler is add event handler.
func AddHandler(name string, handler Handler) error {
	return masterEventMgr.addHandler(masterEventMgr.router, name, handler)
}

// AddHandlerWithResult add event hander with responce callback.
func AddHandlerWithResult(name string, handler HandlerWithResult) error {
	return masterEventMgr.addHandler(masterEventMgr.router, name, handler)
}

// AddGroup is add event group
func AddGroup(name string) (Group, error) {
	return masterEventMgr.addGroup(masterEventMgr.router, name)
}

// Emit is dispatch event
func Emit(name string, argument Value) error {
	return masterEventMgr.emit(name, argument, nil)
}

// EmitWithCallback is dispatch event with callback
func EmitWithCallback(name string, argument Value, callback ResponceCallback) error {
	return masterEventMgr.emit(name, argument, callback)
}

// EmitWithResult is dispatch event by syncronous and return result
func EmitWithResult(name string, argument Value) Result {
	ch := make(chan Result)
	err := EmitWithCallback(name, argument, func(result Result) {
		ch <- result
	})
	if err != nil {
		return NewErrorResult(err)
	}
	return <-ch
}

// SetUnmatchedHandler set handler called at event is unmatched
func SetUnmatchedHandler(handler UnmatchedHandler) {
	masterEventMgr.router.SetUnmatchedItem(handler)
}

// StartEventMgr start event manager
func StartEventMgr() {
	masterEventMgr = &eventMgr{
		router:    NewRouter(),
		eventChan: make(chan cmd),
	}
	go masterEventMgr.eventMain()
}

// StopEventMgr stop event manager
func StopEventMgr() {
	masterEventMgr.quit()
	masterEventMgr = nil
}
