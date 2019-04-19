package event

import (
	"errors"
	"fmt"
	"strings"
)

const HostID = "-"

type EventPathItem struct {
	Name  string
	Value string
}

type EventPath []EventPathItem

func (ep EventPath) String() string {
	// for debug
	var s string
	for _, itm := range ep {
		s += "/" + itm.Name + "/" + itm.Value
	}
	return s
}

func eventTargeetToPathStringCore(target EventTarget) []string {
	var result []string
	if target.ParentTarget() != nil {
		result = eventTargeetToPathStringCore(target.ParentTarget())
	}
	if !target.Host().IsSingleton() {
		result = append(result, target.Host().Name(), target.TargetID())
	}
	return result
}

func EventTargetToPathString(target EventTarget) string {
	p := eventTargeetToPathStringCore(target)
	return "/" + strings.Join(p, "/")
}

func toDriverEventPathCore(target EventTarget, params map[string]string) []string {
	var result []string
	if target.ParentTarget() != nil {
		result = toDriverEventPathCore(target.ParentTarget(), params)
	}
	if !target.Host().IsSingleton() {
		result = append(result, target.Host().Name(), ":"+target.Host().Name())
		params[target.Host().Name()] = target.TargetID()
	}
	return result
}

func ToDriverEventPath(target EventTarget) (string, map[string]string) {
	params := make(map[string]string)
	p := toDriverEventPathCore(target, params)
	if len(p) == 0 {
		if !target.Host().IsSingleton() {
			panic("invalid path")
		}
		p = append(p, target.Host().Name())
	}
	return "/" + strings.Join(p, "/"), params
}

func ToEventPath(target EventTarget) EventPath {
	depth := 1
	for t := target; t.ParentTarget() != nil; t = t.ParentTarget() {
		depth++
	}
	// log.Printf("target = %v, depth = %d\n", target, depth)

	p := make(EventPath, depth)
	depth--
	for t := target; t != nil; t = t.ParentTarget() {
		h := t.Host()
		p[depth].Name = h.Name()
		p[depth].Value = t.TargetID()
		depth--
	}
	// log.Printf("ToEventPath: p = %v\n", p)
	return p
}

func StringToEventTarget(path string) (EventTarget, error) {
	paths, err := StringToEventPath(path)
	if err != nil {
		return nil, err
	}
	// log.Printf("paths = %v,  name = %q\n", paths, n)
	_, t, _ := rootHost.Resolve(paths)
	if t == nil {
		return nil, fmt.Errorf("target is not found: %q", path)
	}
	return t, nil
}

func StringToEventPath(path string) (EventPath, error) {
	if path == "" {
		return nil, errors.New("empty path")
	}
	paths := strings.Split(path, "/")
	if paths[0] == "" {
		paths = paths[1:]
	}
	p := make(EventPath, 0, len(paths))
	var ph EventHost = rootHost.Host()
	for len(paths) > 0 {
		name := paths[0]
		// log.Printf("name = %q, paths=%v, ph = %v\n", name, paths, ph)
		var val string
		if ph.Name() != name {
			if ph.IsSingleton() {
				ph = ph.GetChild(name)
				if ph != nil {
					continue
				}
			}
			return nil, fmt.Errorf("invalid path: %q", path)
		}
		if ph.IsSingleton() {
			val = HostID
			paths = paths[1:]
		} else {
			if len(paths) == 1 {
				return nil, fmt.Errorf("invalid path: %q", path)
			}
			val = paths[1]
			paths = paths[2:]
		}
		if len(paths) > 0 {
			ph = ph.GetChild(paths[0])
			// log.Printf("ph => %v, %q", ph, paths[0])
			if ph == nil {
				return nil, fmt.Errorf("invalid path: %q", path)
			}
		}
		p = append(p, EventPathItem{Name: name, Value: val})
	}
	return p, nil
}

type EventTarget interface {
	ParentTarget() EventTarget //TODO: 整理する
	Host() EventHost
	TargetID() string
}

type EventTargetWithLocalHandler interface {
	EventTarget
	GetEventHandler(name string) Handler
}

type EventTargetWithSignal interface {
	GetEventSignal(name string) *Signal
}

type EventTargetWithSlot interface {
	GetEventSlot(name string) *Slot
}

type EventTargetWithScopedNameResolver interface {
	GetTargetByScopedName(scopedName string) (EventTarget, string)
}

type EventHandler interface {
	Emit(e *Event, callback ResponceCallback) error
	IsEnable() bool
}

type eventHandlerHandler struct {
	handler Handler
}

func (ehh *eventHandlerHandler) Emit(e *Event, callback ResponceCallback) error {
	err := ehh.handler(e)
	if callback != nil {
		if err != nil {
			callback(NewErrorResult(err))
		} else {
			callback(NewValueResult(nil))
		}
	}
	return err
}

func (ehh *eventHandlerHandler) IsEnable() bool {
	return true
}

type eventHandlerHandlerWithResult struct {
	handler HandlerWithResult
}

func (ehh *eventHandlerHandlerWithResult) Emit(e *Event, callback ResponceCallback) error {
	ehh.handler(e, callback)
	return nil
}

func (ehh *eventHandlerHandlerWithResult) IsEnable() bool {
	return true
}

type eventHandlerSignal struct {
	signal *Signal
}

func (ehs *eventHandlerSignal) Emit(e *Event, callback ResponceCallback) error {
	err := ehs.signal.Emit(e.Argument)
	if callback != nil {
		if err != nil {
			callback(NewErrorResult(err))
		} else {
			callback(NewValueResult(nil))
		}
	}
	return err
}

func (ehs *eventHandlerSignal) IsEnable() bool {
	return ehs.signal.IsEnabled()
}

type eventHandlerSlot struct {
	slot *Slot
}

func (ehs *eventHandlerSlot) Emit(e *Event, callback ResponceCallback) error {
	err := ehs.slot.emit(e)
	if callback != nil {
		if err != nil {
			callback(NewErrorResult(err))
		} else {
			callback(NewValueResult(NewValue(e.Result)))
		}
	}
	return err
}

func (ehs *eventHandlerSlot) IsEnable() bool {
	return ehs.slot.IsEnabled()
}

type EventHost interface {
	EventTarget
	Owner() EventHost
	Name() string
	Core() *EventHostCore
	//	Init(name string, owner EventHost)
	//	InitSingleton(name string, owner EventHost)
	IsSingleton() bool
	AddHandler(name string, h Handler)
	AddHandlerWithResult(name string, h HandlerWithResult)
	GetTarget(id string, parent EventTarget) EventTarget
	AddChild(child EventHost)
	GetChild(name string) EventHost
	Resolve(path EventPath) (EventHost, EventTarget, map[string]string)
	// Emit(path EventPath, name string, argument Value, callback ResponceCallback) error
	GetHandler(name string) EventHandler
}

type EventHostCore struct {
	owner     EventHost
	name      string
	host      EventHost
	singleton bool
	children  map[string]EventHost
	handlers  map[string]EventHandler
}

//var rootHost *EventHostCore
var rootHost EventHost

func InitHost(host EventHost, name string, owner EventHost) {
	c := host.Core()
	c.name = name
	c.owner = owner
	c.host = host
	if owner != nil {
		owner.AddChild(host)
	} else {
		if rootHost != nil {
			panic("already exist event root")
		}
		rootHost = host
	}
}

func InitSingletonRoot(host EventHost, name string) {
	c := host.Core()
	c.name = name
	c.owner = nil
	c.host = host
	c.singleton = true
	if rootHost != nil {
		panic("already exist event root")
	}
	rootHost = host
}

func (ehc *EventHostCore) Core() *EventHostCore {
	return ehc
}

func (ehc *EventHostCore) IsSingleton() bool {
	return ehc.singleton
}

func (ehc *EventHostCore) ParentTarget() EventTarget {
	return ehc.owner
}

func Emit(target EventTarget, name string, argument Value) error {
	if rootHost == nil {
		return errors.New("event host is not registerd")
	}
	return emitEventCore(target, name, argument, nil)
}

func EmitWithCallback(target EventTarget, name string, argument Value, callback ResponceCallback) error {
	if rootHost == nil {
		return errors.New("event host is not registerd")
	}
	return emitEventCore(target, name, argument, callback)
}

func IsEnableEvent(target EventTarget, name string) (bool, error) {
	eh, _, _, err := findEventHandler(target, name)
	if err != nil {
		return false, err
	}
	return eh.IsEnable(), nil
}

func findEventHandler(target EventTarget, name string) (EventHandler, EventTarget, EventHost, error) {
	// TODO : change path to Target?
	// log.Printf("EventHostCore::findEventHandler(%v, %q)", target, name)
	host := target.Host()
	if host == nil || target == nil {
		return nil, nil, nil, errors.New("target not found in path")
	}
	if tsn, ok := target.(EventTargetWithScopedNameResolver); ok {
		t, n := tsn.GetTargetByScopedName(name)
		if t != nil {
			target = t
			host = t.Host()
			name = n
		}
	}
	// log.Printf("======> target: %v, host: %v\n", target, host)
	if etlh, ok := target.(EventTargetWithLocalHandler); ok {
		if h := etlh.GetEventHandler(name); h != nil {
			return &eventHandlerHandler{handler: h}, target, host, nil
		}
	}
	if ets, ok := target.(EventTargetWithSignal); ok {
		if s := ets.GetEventSignal(name); s != nil {
			return &eventHandlerSignal{signal: s}, target, host, nil
		}
	}
	if ets, ok := target.(EventTargetWithSlot); ok {
		if s := ets.GetEventSlot(name); s != nil {
			return &eventHandlerSlot{slot: s}, target, host, nil
		}
	}
	if h := host.GetHandler(name); h != nil {
		return h, target, host, nil
	}
	return nil, nil, nil, fmt.Errorf("event %q not found in %v (host: %v)", name, target, host)
}

func emitEventCore(target EventTarget, name string, argument Value, callback ResponceCallback) error {
	eh, target, host, err := findEventHandler(target, name)
	if err != nil {
		// log.Printf("*************** err: %v", err)
		if callback != nil {
			callback(NewErrorResult(err))
		}
		return err
	}
	e := &Event{
		Name:     name,
		Argument: argument,
		Target:   target,
		Host:     host,
	}
	return eh.Emit(e, callback)
	// // TODO : change path to Target?
	// //log.Printf("EventHostCore::Emit(%v, %q, %v, %v)", path, name, argument, callback)
	// host := target.Host()
	// if host == nil || target == nil {
	// 	return errors.New("target not found in path")
	// }
	// if tsn, ok := target.(EventTargetWithScopedNameResolver); ok {
	// 	t, n := tsn.GetTargetByScopedName(name)
	// 	if t != nil {
	// 		target = t
	// 		host = t.Host()
	// 		name = n
	// 	}
	// }
	// //log.Printf("======> target: %v, host: %v, params: %v\n", target, host, params)
	// e := &Event{
	// 	Name:     name,
	// 	Argument: argument,
	// 	Target:   target,
	// 	Host:     host,
	// }
	// if etlh, ok := target.(EventTargetWithLocalHandler); ok {
	// 	if h := etlh.GetEventHandler(name); h != nil {
	// 		if h != nil {
	// 			h(e)
	// 			if callback != nil {
	// 				r := e.Result
	// 				if e, ok := r.(error); ok {
	// 					callback(NewErrorResult(e))
	// 					return e
	// 				}
	// 				callback(NewValueResult(NewValue(r)))
	// 			}
	// 			return nil
	// 		}
	// 	}
	// }
	// if ets, ok := target.(EventTargetWithSignal); ok {
	// 	if s := ets.GetEventSignal(name); s != nil {
	// 		err := s.Emit(argument)
	// 		if callback != nil {
	// 			callback(NewErrorResult(err))
	// 		}
	// 		return err
	// 	}
	// }
	// if ets, ok := target.(EventTargetWithSlot); ok {
	// 	if s := ets.GetEventSlot(name); s != nil {
	// 		err := s.emit(e)
	// 		if callback != nil {
	// 			callback(NewErrorResult(err))
	// 		}
	// 		return err
	// 	}
	// }
	// if h := host.GetHandler(name); h != nil {
	// 	return h.Emit(e, callback)
	// }
	// return fmt.Errorf("event %q not found in %v (host: %v)", name, target, host)
}

func EmitWithResult(target EventTarget, name string, argument Value) Result {
	rc := make(chan Result)
	err := EmitWithCallback(target, name, argument, func(e Result) {
		rc <- e
	})
	if err != nil {
		return NewErrorResult(err)
	}
	return <-rc
}

func (ehc *EventHostCore) Host() EventHost {
	return ehc.host
}

func (ehc *EventHostCore) TargetID() string {
	return HostID
}

func (ehc *EventHostCore) Owner() EventHost {
	return ehc.owner
}

func (ehc *EventHostCore) Name() string {
	return ehc.name
}

func (ehc *EventHostCore) AddHandler(name string, h Handler) {
	if ehc.handlers == nil {
		ehc.handlers = make(map[string]EventHandler)
	}
	// log.Printf("EventHostCore::AddHandler(%v, %q, %v)", ehc, name, h)
	ehc.handlers[name] = &eventHandlerHandler{handler: h}
}

func (ehc *EventHostCore) AddHandlerWithResult(name string, h HandlerWithResult) {
	if ehc.handlers == nil {
		ehc.handlers = make(map[string]EventHandler)
	}
	ehc.handlers[name] = &eventHandlerHandlerWithResult{handler: h}
}

func (ehc *EventHostCore) AddChild(child EventHost) {
	if ehc.children == nil {
		ehc.children = make(map[string]EventHost)
	}
	ehc.children[child.Name()] = child
}

func (ehc *EventHostCore) GetChild(name string) EventHost {
	if ehc.children == nil {
		return nil
	}
	return ehc.children[name]
}

func (ehc *EventHostCore) Resolve(path EventPath) (EventHost, EventTarget, map[string]string) {
	//TODO: remove EventHost from return value?
	ph := ehc.Host()
	var target EventTarget
	var host EventHost
	params := make(map[string]string)
	// log.Printf("### path: %v", path)
	for len(path) > 0 {
		// log.Printf("#### Path: %q", path)
		n := path[0].Name
		v := path[0].Value

		// log.Printf("##### Name: %q, v: %q", n, v)
		if ph.Name() != n {
			if ph.IsSingleton() {
				ph = ph.GetChild(n)
				if ph != nil {
					continue
				}
			}
			// log.Printf("### => not found(n = %q, ph = %v)\n", n, ph)
			return nil, nil, nil
		}
		target = ph.GetTarget(v, target)
		host = ph
		params[n] = v
		path = path[1:]
		if len(path) > 0 {
			ph = ph.GetChild(path[0].Name)
		}
	}
	// log.Printf("### => %v, %v, %v", host, target, params)
	return host, target, params
}

// func (ehc *EventHostCore) Emit(path EventPath, name string, argument Value, callback ResponceCallback) error {
// 	// TODO : change path to Target?
// 	//log.Printf("EventHostCore::Emit(%v, %q, %v, %v)", path, name, argument, callback)
// 	host, target, _ := ehc.Resolve(path) // params 要らんかも
// 	if host == nil || target == nil {
// 		return errors.New("target not found in path")
// 	}
// 	if tsn, ok := target.(EventTargetWithScopedNameResolver); ok {
// 		t, n := tsn.GetTargetByScopedName(name)
// 		if t != nil {
// 			target = t
// 			host = t.Host()
// 			name = n
// 		}
// 	}
// 	//log.Printf("======> target: %v, host: %v, params: %v\n", target, host, params)
// 	e := &Event{
// 		Name:     name,
// 		Argument: argument,
// 		Target:   target,
// 		Host:     host,
// 	}
// 	if etlh, ok := target.(EventTargetWithLocalHandler); ok {
// 		if h := etlh.GetEventHandler(name); h != nil {
// 			if h != nil {
// 				h(e)
// 				if callback != nil {
// 					r := e.Result
// 					if e, ok := r.(error); ok {
// 						callback(NewErrorResult(e))
// 						return e
// 					}
// 					callback(NewValueResult(NewValue(r)))
// 				}
// 				return nil
// 			}
// 		}
// 	}
// 	if ets, ok := target.(EventTargetWithSignal); ok {
// 		if s := ets.GetEventSignal(name); s != nil {
// 			err := s.Emit(argument)
// 			if callback != nil {
// 				callback(NewErrorResult(err))
// 			}
// 			return err
// 		}
// 	}
// 	if ets, ok := target.(EventTargetWithSlot); ok {
// 		if s := ets.GetEventSlot(name); s != nil {
// 			err := s.emit(e)
// 			if callback != nil {
// 				callback(NewErrorResult(err))
// 			}
// 			return err
// 		}
// 	}
// 	if h := host.GetHandler(name); h != nil {
// 		return h.Emit(e, callback)
// 	}
// 	return fmt.Errorf("event %q not found in %v (host: %v)", name, target, host)
// }

func (ehc *EventHostCore) GetHandler(name string) EventHandler {
	if h, ok := ehc.handlers[name]; ok {
		return h
	}
	return nil
}
