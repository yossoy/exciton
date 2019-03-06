package windows

/*
#include "driver.h"
#include "menu.h"
*/
import "C"

import (
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
)

func initializeMenu(serializer driver.DriverEventSerializer) error {
	menu.MenuClass.AddHandlerWithResult("new", func(e *event.Event, callback event.ResponceCallback) {
		serializer.RelayEventWithResult(e, callback)
	})
	menu.MenuClass.AddHandlerWithResult("updateDiffSetHandler", func(e *event.Event, callback event.ResponceCallback) {
		serializer.RelayEventWithResult(e, callback)
	})
	menu.MenuClass.AddHandler("setApplicationMenu", func(e *event.Event) error {
		serializer.RelayEvent(e)
		return nil
	})
	menu.MenuClass.AddHandler("popupContextMenu", func(e *event.Event) error {
		serializer.RelayEvent(e)
		return nil
	})

	C.Menu_Init()
	return nil
}
