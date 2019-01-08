package exciton

import (
	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
)

type InternalInitFunc func(app *app.App, info *app.StartupInfo) error

func Init(info *app.StartupInfo, initFunc InternalInitFunc) error {
	rootGroup := info.AppEventRoot
	if info.OnAppQuit != nil {
		rootGroup.AddHandler("/app/finalize", func(e *event.Event) {
			info.OnAppQuit()
		})
	}
	if err := rootGroup.AddHandler("/app/finalizedWindow", func(e *event.Event) {
		app, err := app.GetAppFromEvent(e)
		if err != nil {
			panic(err) //TODO: error handling
		}
		var win *window.Window
		if err := e.Argument.Decode(&win); err != nil {
			panic(err)
		}
		if win == app.MainWindow {
			app.MainWindow = nil
		}
	}); err != nil {
		return err
	}
	if err := rootGroup.AddHandler("/app/init", func(e *event.Event) {
		app, err := app.GetAppFromEvent(e)
		if err != nil {
			panic(err) //TODO: error handling
		}
		err = initFunc(app, info)
		if err != nil {
			panic(err) //TODO: error handlig
		}
	}); err != nil {
		return err
	}
	if err := window.InitWindows(&info.StartupInfo); err != nil {
		return err
	}
	if err := menu.InitMenus(rootGroup); err != nil {
		return err
	}
	if err := markup.InitEvents(rootGroup); err != nil {
		return err
	}
	return nil
}
