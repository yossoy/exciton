package menu

import (
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/geom"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/window"
)

func ContextMenu(m ...markup.MarkupOrChild) *markup.RenderResult {
	return markup.Tag("menu", m...)
}

type popupContextMenuArg struct {
	Position geom.Point `json:"position"`
	WindowID string     `json:"windowId"`
}

func (m *MenuInstance) Popup(mousePt geom.Point, parent *window.Window) error {
	arg := popupContextMenuArg{
		Position: mousePt,
		WindowID: parent.ID,
	}
	return event.Emit("/menu/"+m.uuid+"/popupContextMenu", event.NewValue(&arg))
}
