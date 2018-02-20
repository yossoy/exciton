package markup

import (
	"github.com/yossoy/exciton/event"
	"github.com/pkg/errors"
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
