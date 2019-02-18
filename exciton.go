package exciton

import (
	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/window"
)

type InternalInitFunc func(app *app.App, info *app.StartupInfo) error

func Init(info *app.StartupInfo, initFunc InternalInitFunc) error {
	app.AppClass.AddHandler("finalize", func(e *event.Event) {
		if info.OnAppQuit != nil {
			info.OnAppQuit()
		}
	})
	app.AppClass.AddHandler("finalizedWindow", func(e *event.Event) {
		a := app.GetAppFromEvent(e)
		if a == nil {
			panic("invalid target")
		}
		var win *window.Window
		if err := e.Argument.Decode(&win); err != nil {
			panic(err)
		}
		if win == a.MainWindow {
			a.MainWindow = nil
		}
	})
	app.AppClass.AddHandler("init", func(e *event.Event) {
		a := app.GetAppFromEvent(e)
		if a == nil {
			panic("invalid target")
		}
		err := initFunc(a, info)
		if err != nil {
			panic(err)
		}
	})
	if err := window.InitWindows(&info.StartupInfo); err != nil {
		return err
	}
	return nil
}
