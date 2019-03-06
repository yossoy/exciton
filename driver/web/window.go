package web

import (
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/window"
)

func initializeWindow(serializer driver.DriverEventSerializer) error {
	window.WindowClass.AddHandlerWithResult("new", func(e *event.Event, callback event.ResponceCallback) {
		serializer.RelayEventWithResult(e, callback)
	})
	window.WindowClass.AddHandler("requestAnimationFrame", func(e *event.Event) error {
		serializer.RelayEvent(e)
		return nil
	})
	window.WindowClass.AddHandler("updateDiffSetHandler", func(e *event.Event) error {
		serializer.RelayEvent(e)
		return nil
	})
	window.WindowClass.AddHandlerWithResult("browserSync", func(e *event.Event, callback event.ResponceCallback) {
		serializer.RelayEventWithResult(e, callback)
	})
	window.WindowClass.AddHandler("browserAsync", func(e *event.Event) error {
		serializer.RelayEvent(e)
		return nil
	})
	window.WindowClass.AddHandler("redirectTo", func(e *event.Event) error {
		serializer.RelayEvent(e)
		return nil
	})
	return nil
}
