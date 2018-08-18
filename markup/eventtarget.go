package markup

import (
	"github.com/yossoy/exciton/internal/object"
)

type EventTarget struct {
	AppID     string `json:"appId,omitempty"`
	WindowID  string `json:"windowId,omitempty"`
	MenuID    string `json:"menuId,omitempty"`
	ElementID string `json:"elementId,omitempty"`
}

func (et *EventTarget) Builder() Builder {
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

func (et *EventTarget) Node() Node {
	if et == nil {
		return nil
	}
	b := et.Builder()
	if b == nil {
		return nil
	}
	itm := b.(*builder).elements.Get(et.ElementID)
	if n, ok := itm.(*node); ok {
		return n
	}
	return nil
}

func (et *EventTarget) HostComponent() Component {
	b := et.Builder()
	if b == nil {
		return nil
	}
	bb := b.(*builder)
	itm := bb.elements.Get(et.ElementID)
	n, ok := itm.(*node)
	if !ok {
		return nil
	}
	for n != nil {
		if n == bb.rootNode {
			break
		}
		if n.component != nil {
			return n.component
		}
		n = n.parent
	}
	return nil
}
