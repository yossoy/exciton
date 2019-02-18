package web

import (
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/menu"
)

func initializeMenu(serializer driver.DriverEventSerializer) error {
	menu.MenuClass.AddHandlerWithResult("new", func(e *event.Event, callback event.ResponceCallback) {
		serializer.RelayEventWithResult(e, callback)
	})
	menu.MenuClass.AddHandlerWithResult("updateDiffSetHandler", func(e *event.Event, callback event.ResponceCallback) {
		serializer.RelayEventWithResult(e, callback)
	})
	menu.MenuClass.AddHandler("setApplicationMenu", func(e *event.Event) {
		serializer.RelayEvent(e)
	})
	menu.MenuClass.AddHandler("popupContextMenu", func(e *event.Event) {
		serializer.RelayEvent(e)
	})

	return nil
}
