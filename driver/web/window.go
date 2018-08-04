package web

import (
	"github.com/yossoy/exciton/event"
)

func initializeWindow(gg event.Group) error {
	g, err := gg.AddGroup("/window/:id")
	if err != nil {
		return err
	}
	g.AddHandlerWithResult("/new", func(e *event.Event, callback event.ResponceCallback) {
		platform.relayEventWithResultToNative(e, callback)
		driverLogDebug("window/new: %v", e.Name)
		//callback(event.NewValueResult(event.NewValue(0)))
	})
	g.AddHandler("/requestAnimationFrame", func(e *event.Event) {
		platform.relayEventToNative(e)
		driverLogDebug("window/requestAnimationFrame")
		//TODO:
	})
	g.AddHandler("/updateDiffSetHandler", func(e *event.Event) {
		platform.relayEventToNative(e)
		driverLogDebug("window/updateDiffSetHandler")
		//TODO:
	})
	g.AddHandlerWithResult("/browserSync", func(e *event.Event, callback event.ResponceCallback) {
		platform.relayEventWithResultToNative(e, callback)
		driverLogDebug("window/browserSync")
		//TODO:
	})
	g.AddHandler("/redirectTo", func(e *event.Event) {
		platform.relayEventToNative(e)
		driverLogDebug("window/redirectTo")
		//TODO:
	})
	return nil
}
