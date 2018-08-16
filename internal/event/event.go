package event

import (
	"errors"

	"github.com/yossoy/exciton/event"
)

// UnmatchedHandler is called at unmatched path
type UnmatchedHandler func(path string, parameter map[string]string, callback event.ResponceCallback)

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
	respCallback     event.ResponceCallback
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

type eventGroup struct {
	subRouter *Router
}

type rootGroup struct {
}

func (rg *rootGroup) AddHandler(name string, handler event.Handler) error {
	return AddHandler(name, handler)
}
func (rg *rootGroup) AddHandlerWithResult(name string, handler event.HandlerWithResult) error {
	return AddHandlerWithResult(name, handler)
}
func (rg *rootGroup) AddGroup(name string) (event.Group, error) {
	return AddGroup(name)
}
func (rg *rootGroup) SetUnmatchedHandler(handler UnmatchedHandler) {
	SetUnmatchedHandler(handler)
}

func RootGroup() event.Group {
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

func (em *eventMgr) addGroup(tr *Router, name string) (event.Group, error) {
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

func (em *eventMgr) emit(name string, argument event.Value, respCallback event.ResponceCallback) error {
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
					var arg event.Value
					if c.argument != nil {
						arg = c.argument.(event.Value)
					}
					e := &event.Event{
						Name:     evn,
						Argument: arg,
						Params:   params,
					}
					if u, ok := r.(UnmatchedRouteItem); ok {
						uh := u.Item().(UnmatchedHandler)
						go uh(u.PathSegments(), params, c.respCallback)
					} else if h, ok := r.Item().(event.Handler); ok {
						if c.respCallback != nil {
							c.resultChan <- errors.New("call EmitWithCallback by non-result event")
						} else {
							go h(e)
							c.resultChan <- nil
						}
					} else if h, ok := r.Item().(event.HandlerWithResult); ok {
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

func (eg *eventGroup) AddHandler(name string, handler event.Handler) error {
	return masterEventMgr.addHandler(eg.subRouter, name, handler)
}

func (eg *eventGroup) AddHandlerWithResult(name string, handler event.HandlerWithResult) error {
	return masterEventMgr.addHandler(eg.subRouter, name, handler)
}

func (eg *eventGroup) AddGroup(name string) (event.Group, error) {
	return masterEventMgr.addGroup(eg.subRouter, name)
}

func (eg *eventGroup) SetUnmatchedHandler(handler UnmatchedHandler) {
	eg.subRouter.SetUnmatchedItem(handler)
}

var (
	masterEventMgr *eventMgr
)

// AddHandler is add event handler.
func AddHandler(name string, handler event.Handler) error {
	return masterEventMgr.addHandler(masterEventMgr.router, name, handler)
}

// AddHandlerWithResult add event hander with responce callback.
func AddHandlerWithResult(name string, handler event.HandlerWithResult) error {
	return masterEventMgr.addHandler(masterEventMgr.router, name, handler)
}

// AddGroup is add event group
func AddGroup(name string) (event.Group, error) {
	return masterEventMgr.addGroup(masterEventMgr.router, name)
}

// Emit is dispatch event
func Emit(name string, argument event.Value) error {
	return masterEventMgr.emit(name, argument, nil)
}

// EmitWithCallback is dispatch event with callback
func EmitWithCallback(name string, argument event.Value, callback event.ResponceCallback) error {
	return masterEventMgr.emit(name, argument, callback)
}

// EmitWithResult is dispatch event by syncronous and return result
func EmitWithResult(name string, argument event.Value) event.Result {
	ch := make(chan event.Result)
	err := EmitWithCallback(name, argument, func(result event.Result) {
		ch <- result
	})
	if err != nil {
		return event.NewErrorResult(err)
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
