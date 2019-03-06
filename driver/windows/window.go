package windows

/*
#cgo CFLAGS: -DUNICODE -D_UNICODE -DWIN32 -D_WINDOWS -D_WIN32_IE=0x0A00 -D_WIN32_WINNT=0x0A00
#cgo CXXFLAGS: -DUNICODE -D_UNICODE -DWIN32 -D_WINDOWS -D_WIN32_IE=0x0A00 -D_WIN32_WINNT=0x0A00 -std=c++17
#cgo LDFLAGS: -static-libgcc -static-libstdc++ -lgdi32 -lole32 -lcomctl32 -loleaut32 -luuid -lurlmon -lwininet -lmshtml -lversion -lshlwapi -Wl,-Bstatic -lstdc++ -lpthread

#include "window.h"
*/
import "C"

import (
	"github.com/yossoy/exciton/window"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
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
	C.Window_Init()
	return nil
}
