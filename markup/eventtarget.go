package markup

import (
	"github.com/pkg/errors"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/internal/object"
)

type EventTarget struct {
	AppID     string `json:"appId,omitempty"`
	WindowID  string `json:"windowId,omitempty"`
	MenuID    string `json:"menuId,omitempty"`
	ElementID string `json:"elementId,omitempty"`
}

type browserCommand struct {
	Command  string      `json:"cmd"`
	Target   interface{} `json:"target"`
	Argument interface{} `json:"argument"`
}

func (et *EventTarget) Builder() *Builder {
	var buildable Buildable
	switch {
	case et.WindowID != "":
		buildable = object.Windows.Get(et.WindowID).(Buildable)
	case et.MenuID != "":
		buildable = object.Menus.Get(et.MenuID).(Buildable)
	default:
		return nil
	}
	return buildable.Builder()
}

func (et *EventTarget) eventRoot() string {
	var buildable Buildable
	switch {
	case et.WindowID != "":
		buildable = object.Windows.Get(et.WindowID).(Buildable)
	case et.MenuID != "":
		buildable = object.Menus.Get(et.MenuID).(Buildable)
	default:
		return ""
	}
	return buildable.EventRoot()
}

func (et *EventTarget) Node() *node {
	builder := et.Builder()
	if builder == nil {
		return nil
	}
	itm := builder.elements.Get(et.ElementID)
	if n, ok := itm.(*node); ok {
		return n
	}
	return nil
}

func (et *EventTarget) HostComponent() Component {
	builder := et.Builder()
	if builder == nil {
		return nil
	}
	itm := builder.elements.Get(et.ElementID)
	n, ok := itm.(*node)
	if !ok {
		return nil
	}
	for n != nil {
		if n == builder.rootNode {
			break
		}
		if n.component != nil {
			return n.component
		}
		n = n.parent
	}
	return nil
}

func (et *EventTarget) GetProperty(name string) (interface{}, error) {
	if et.ElementID == "" {
		// target is window
		panic("not implement yet")
	}
	arg := &browserCommand{
		Command:  "getProp",
		Target:   et.ElementID,
		Argument: name,
	}
	eventRoot := et.eventRoot()
	var result event.Result
	if et.WindowID != "" {
		result = event.EmitWithResult(eventRoot+"/window/"+et.WindowID+"/browserSync", event.NewValue(arg))
	} else if et.MenuID != "" {
		result = event.EmitWithResult(eventRoot+"/menu"+et.MenuID+"/browserSync", event.NewValue(arg))
	}
	if result.Error() != nil {
		return nil, errors.Wrap(result.Error(), "EmitEventWithResult fail:")
	}
	var ret interface{}
	if err := result.Value().Decode(&ret); err != nil {
		return nil, errors.Wrap(err, "Value.Decode() fail.:")
	}
	return ret, nil
}
