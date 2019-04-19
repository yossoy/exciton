package menu

import (
	"fmt"

	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/geom"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/window"
)

func contextMenu(m ...markup.MarkupOrChild) markup.RenderResult {
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
	return event.Emit(m, "popupContextMenu", event.NewValue(&arg))
}

func PopupMenu(menu MenuTemplate, mousePt geom.Point, w *window.Window) error {
	mi, err := newPopupMenu(w.Owner(), w, menu)
	if err != nil {
		return err
	}
	return mi.Popup(mousePt, w)
}

func PopupMenuOnComponent(menu MenuTemplate, mousePt geom.Point, c markup.Component) error {
	// m, err := toPopupMenuSub(menu)
	// if err != nil {
	// 	return err
	// }
	b := c.Builder().Buildable()
	w, ok := b.(*window.Window)
	if !ok {
		return fmt.Errorf("invalid target")
	}

	mi, err := newPopupMenu(w.Owner(), c, menu)
	if err != nil {
		return err
	}
	return mi.Popup(mousePt, w)
}
