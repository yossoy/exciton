package windows

/*
#include "driver.h"
#include "menu.h"
*/
import "C"

import (
	"github.com/yossoy/exciton/event"
)

func initializeMenu() error {
	g, err := event.AddGroup("/menu/:id")
	if err != nil {
		return err
	}
	g.AddHandlerWithResult("/new", func(e *event.Event, callback event.ResponceCallback) {
		platform.relayEventWithResultToNative(e, callback)
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

	C.Menu_Init()
	return nil
}
