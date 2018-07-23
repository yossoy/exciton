package mac

/*
#include "driver.h"
#include "window.h"
*/
import "C"
import "github.com/yossoy/exciton/event"

func initializeWindow() error {
	g, err := event.AddGroup("/window/:id")
	if err != nil {
		return err
	}
	g.AddHandlerWithResult("/new", func(e *event.Event, callback event.ResponceCallback) {
		platform.relayEventWithResultToNative(e, callback)
	})
	g.AddHandler("/requestAnimationFrame", func(e *event.Event) {
		platform.relayEventToNative(e)
	})
	g.AddHandler("/updateDiffSetHandler", func(e *event.Event) {
		platform.relayEventToNative(e)
	})
	g.AddHandlerWithResult("/browserSync", func(e *event.Event, callback event.ResponceCallback) {
		platform.relayEventWithResultToNative(e, callback)
	})
	g.AddHandler("/redirectTo", func(e *event.Event) {
		platform.relayEventToNative(e)
	})
	C.Window_Init()
	return nil
}
