package markup

import (
	"github.com/pkg/errors"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/internal/object"
)

type EventTarget struct {
	WindowID  string `json:"windowId,omitempty"`
	MenuID    string `json:"menuId,omitempty"`
	ElementID string `json:"elementId,omitempty"`
}

type browserCommand struct {
	Command   string `json:"cmd"`
	ElementID string `json:"elemId"`
	Property  string `json:"propName"`
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

func (et *EventTarget) GetProperty(name string) (interface{}, error) {
	if et.ElementID == "" {
		// target is window
		panic("not implement yet")
	}
	arg := &browserCommand{
		"getProp",
		et.ElementID,
		name,
	}
	var result event.Result
	if et.WindowID != "" {
		result = event.EmitWithResult("/window/"+et.WindowID+"/browserSync", event.NewValue(arg))
	} else if et.MenuID != "" {
		result = event.EmitWithResult("/menu"+et.MenuID+"/browserSync", event.NewValue(arg))
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
