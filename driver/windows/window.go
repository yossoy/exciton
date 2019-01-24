package windows

/*
#cgo CFLAGS: -DUNICODE -D_UNICODE -DWIN32 -D_WINDOWS -D_WIN32_IE=0x0A00 -D_WIN32_WINNT=0x0A00
#cgo CXXFLAGS: -DUNICODE -D_UNICODE -DWIN32 -D_WINDOWS -D_WIN32_IE=0x0A00 -D_WIN32_WINNT=0x0A00 -std=c++17
#cgo LDFLAGS: -static-libgcc -static-libstdc++ -lgdi32 -lole32 -lcomctl32 -loleaut32 -luuid -lurlmon -lwininet -lmshtml -lversion -lshlwapi -Wl,-Bstatic -lstdc++ -lpthread

#include "window.h"
*/
import "C"

import (
	"github.com/yossoy/exciton/event"
	ievent "github.com/yossoy/exciton/internal/event"
)

func initializeWindow() error {
	g, err := ievent.AddGroup("/window/:id")
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
	g.AddHandler("/browserAsync", func(e *event.Event) {
		platform.relayEventToNative(e)
	})
	g.AddHandler("/redirectTo", func(e *event.Event) {
		platform.relayEventToNative(e)
	})
	C.Window_Init()
	return nil
}
