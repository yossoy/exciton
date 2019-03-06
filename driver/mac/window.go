package mac

/*
#include "driver.h"
#include "window.h"
*/
import "C"
import (
	"github.com/yossoy/exciton/window"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
)

func initializeWindow(serializer driver.DriverEventSerializer) error {
	window.WindowClass.AddHandlerWithResult("new", func (e *event.Event, callback event.ResponceCallback) {
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
	C.Window_Init()
	return nil
}
