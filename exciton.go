package exciton

import (
	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/window"
)

type InternalInitFunc func(app *app.App, info *app.StartupInfo) error

func Init(info *app.StartupInfo, initFunc InternalInitFunc) error {
	app.AppClass.AddHandler("finalize", func(e *event.Event) error {
		if info.OnAppQuit != nil {
			info.OnAppQuit()
		}
		return nil
	})
	app.AppClass.AddHandler("finalizedWindow", func(e *event.Event) error {
		a, err := app.GetAppFromEvent(e)
		if err != nil {
			return err
		}
		var win *window.Window
		if err := e.Argument.Decode(&win); err != nil {
			return err
		}
		if win == a.MainWindow {
			a.MainWindow = nil
		}
		return nil
	})
	app.AppClass.AddHandler("init", func(e *event.Event) error {
		a, err := app.GetAppFromEvent(e)
		if err != nil {
			return err
		}
		err = initFunc(a, info)
		if err != nil {
			return err
		}
		return nil
	})
	if err := window.InitWindows(&info.StartupInfo); err != nil {
		return err
	}
	return nil
}
