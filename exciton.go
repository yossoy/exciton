package exciton

import (
	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
)

type InternalInitFunc func(app *app.App, info *app.StartupInfo) error

func Init(rootGroup event.Group, info *app.StartupInfo, initFunc InternalInitFunc) error {
	if info.OnAppQuit != nil {
		rootGroup.AddHandler("/app/finalize", func(e *event.Event) {
			info.OnAppQuit()
		})
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
	if err := window.InitWindows(rootGroup, &info.StartupInfo); err != nil {
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
