package menu

import (
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/markup"
)

func ApplicationMenu(m ...markup.MarkupOrChild) markup.RenderResult {
	return markup.Tag("menu", m...)
}

func SetApplicationMenu(m *MenuInstance) error {
	return event.Emit("/menu/"+m.uuid+"/setApplicationMenu", event.NewValue(nil))
}
