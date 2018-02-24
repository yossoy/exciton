package windows

/*
#cgo CFLAGS: -DUNICODE -D_UNICODE -DWIN32 -D_WINDOWS -D_WIN32_IE=0x0A00 -D_WIN32_WINNT=0x0A00
#cgo CXXFLAGS: -DUNICODE -D_UNICODE -DWIN32 -D_WINDOWS -D_WIN32_IE=0x0A00 -D_WIN32_WINNT=0x0A00 -std=c++17
#cgo LDFLAGS: -static-libgcc -static-libstdc++ -lgdi32 -lole32 -lcomctl32 -loleaut32 -luuid -lurlmon -lwininet -lmshtml -lmshtml

#include "window.h"
*/
import "C"

import (
	"github.com/yossoy/exciton/event"
)

func initializeWindow() error {
	g, err := event.AddGroup("/window/:id")
	if err != nil {
		return err
	}
	err = g.AddHandlerWithResult("/new", func(e *event.Event, callback event.ResponceCallback) {
		platform.relayEventWithResultToNative(e, callback)
	})
	if err != nil {
		return err
	}
	g.AddHandler("/requestAnimationFrame", func(e *event.Event) {
		platform.relayEventToNative(e)
	})
	g.AddHandler("/updateDiffSetHandler", func(e *event.Event) {
		platform.relayEventToNative(e)
	})
	g.AddHandlerWithResult("/browserSync", func(e *event.Event, callback event.ResponceCallback) {
		platform.relayEventWithResultToNative(e, callback)
	})
	C.Window_Init()
	return nil
}
