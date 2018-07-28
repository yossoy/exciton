package exciton

import (
	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
)

type InternalInitFunc func(info *app.StartupInfo) error

func Init(info *app.StartupInfo, initFunc InternalInitFunc) error {
	if info.OnAppQuit != nil {
		event.AddHandler("/app/finalize", func(e *event.Event) {
			info.OnAppQuit()
		})
	}
	if err := event.AddHandler("/app/init", func(e *event.Event) {
		err := initFunc(info)
		if err != nil {
			panic(err) //TODO: error handlig
		}
	}); err != nil {
		return err
	}
	if err := window.InitWindows(&info.StartupInfo); err != nil {
		return err
	}
	if err := menu.InitMenus(); err != nil {
		return err
	}
	if err := markup.InitEvents(); err != nil {
		return err
	}
	return nil
}
