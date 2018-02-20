package mac

/*
#include "driver.h"
#include "dialog.h"
*/
import "C"
import (
	"github.com/yossoy/exciton/event"
)

func initializeDialog() error {
	g, err := event.AddGroup("/dialog/:id")
	if err != nil {
		return err
	}
	g.AddHandlerWithResult("/showMessageBox", func(e *event.Event, callback event.ResponceCallback) {
		platform.relayEventWithResultToNative(e, callback)
	})
	g.AddHandlerWithResult("/showOpenDialog", func(e *event.Event, callback event.ResponceCallback) {
		platform.relayEventWithResultToNative(e, callback)
	})
	g.AddHandlerWithResult("/showSaveDialog", func(e *event.Event, callback event.ResponceCallback) {
		platform.relayEventWithResultToNative(e, callback)
	})

	C.Dialog_Init()
	return nil
}
