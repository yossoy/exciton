package web

import (
	"github.com/yossoy/exciton/event"
)

func initializeMenu(gg event.Group) error {
	g, err := gg.AddGroup("/menu/:id")
	if err != nil {
		return err
	}
	g.AddHandlerWithResult("/new", func(e *event.Event, callback event.ResponceCallback) {
		platform.relayEventWithResultToNative(e, callback)
		driverLogDebug("window/new: %v", e.Name)
	})
	g.AddHandlerWithResult("/updateDiffSetHandler", func(e *event.Event, callback event.ResponceCallback) {
		platform.relayEventWithResultToNative(e, callback)
	})
	g.AddHandler("/setApplicationMenu", func(e *event.Event) {
		platform.relayEventToNative(e)
	})
	g.AddHandler("/popupContextMenu", func(e *event.Event) {
		platform.relayEventToNative(e)
	})

	return nil
}
