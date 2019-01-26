package markup

import (
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/internal/object"
	"github.com/yossoy/exciton/log"
)

type EventListener interface {
	Markup
	PreventDefault() EventListener
	StopPropagation() EventListener
}

type eventListener struct {
	EventListener
	Name                string
	Listener            event.Handler
	clientScriptPrefix  string
	scriptHandlerName   string
	scriptArguments     []interface{}
	callPreventDefault  bool
	callStopPropagation bool
	//TODO: binding position
}

func NewEventListener(name string, listener event.Handler) *eventListener {
	return &eventListener{
		Name:     name,
		Listener: listener,
	}
}

func NewClientEventListener(name string, scriptPrefix string, handlerName string, arguments []interface{}) *eventListener {
	return &eventListener{
		Name:               name,
		clientScriptPrefix: scriptPrefix,
		scriptHandlerName:  handlerName,
		scriptArguments:    arguments,
	}
}

func (l *eventListener) isMarkup()        {}
func (l *eventListener) isMarkupOrChild() {}

// PreventDefault prevents the default behavior of the event from occurring.
//
// See https://developer.mozilla.org/en-US/docs/Web/API/Event/preventDefault.
func (l *eventListener) PreventDefault() EventListener {
	//TODO: if eventListener is client script, PreventDefault does not affect.
	l.callPreventDefault = true
	return l
}

// StopPropagation prevents further propagation of the current event in the
// capturing and bubbling phases.
//
// See https://developer.mozilla.org/en-US/docs/Web/API/Event/stopPropagation.
func (l *eventListener) StopPropagation() EventListener {
	//TODO: if eventListener is client script, StopPropagation does not affect.
	l.callStopPropagation = true
	return l
}

// Apply implements the Applyer interface.
func (l *eventListener) applyToNode(b Builder, n Node, on Node) {
	bb := b.(*builder)
	nn := n.(*node)
	onn := on.(*node)
	match := false
	exist := false
	if onn.eventListeners != nil {
		if oe, ok := onn.eventListeners[l.Name]; ok {
			if l.callPreventDefault == oe.callPreventDefault && l.callStopPropagation == oe.callStopPropagation {
				if l.Listener != nil && oe.Listener != nil {
					// There is not need to compare Listener,
					//   	 because Listner is used only in go region and is not used in js region.
					match = true
				} else if l.Listener == nil && oe.Listener == nil {
					match = (l.clientScriptPrefix == oe.clientScriptPrefix) && (l.scriptHandlerName == oe.scriptHandlerName)
				}
			}
			exist = true
			delete(onn.eventListeners, l.Name)
		}
	}
	if nn.eventListeners == nil {
		nn.eventListeners = make(map[string]*eventListener)
	}
	nn.eventListeners[l.Name] = l
	if !match {
		if exist {
			bb.diffSet.RemoveEventListener(nn, l.Name)
		}
		if l.Listener != nil {
			bb.diffSet.AddEventListener(nn, l.Name, nn.uuid, l.callPreventDefault, l.callStopPropagation)
		} else {
			bb.diffSet.AddClientEvent(nn, l.Name, nn.uuid, l.clientScriptPrefix, l.scriptHandlerName, l.scriptArguments)
		}
	}
}

func InitEvents(appg event.Group) error {
	err := appg.AddHandler("/:evtroot/:id/html/:html/:event", func(e *event.Event) {
		//id := e.Params["id"]
		eventRoot := e.Params["evtroot"]
		id := e.Params["id"]
		html := e.Params["html"]
		event := e.Params["event"]
		//TODO: cleaup code!
		var buildable Buildable
		switch eventRoot {
		case "window":
			buildable = object.Windows.Get(id).(Buildable)
		case "menu":
			buildable = object.Menus.Get(id).(Buildable)
		default:
			panic("invalid html event path:/" + eventRoot + "/...")
		}
		itm := buildable.Builder().(*builder).elements.Get(html)
		if itm == nil {
			log.PrintError("obj not found: %q", html)
			return
		}
		obj, ok := itm.(*node)
		if !ok {
			panic("obj not found(invalid sequecne): " + html)
		}
		l, ok := obj.eventListeners[event]
		if !ok {
			log.PrintError("event is not registerd: %q", event)
			return
		}
		l.Listener(e)
	})
	return err
}
