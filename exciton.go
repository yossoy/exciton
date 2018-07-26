package exciton

import (
	"errors"

	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/internal/menu"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/window"
)

// RunCallback is called at ready application

type StartupFunc = driver.StartupFunc

type StartupInfo = driver.StartupInfo

func Init(info *driver.StartupInfo) error {
	//event.StartEventMgr()
	if info.OnAppStart == nil {
		return errors.New("Need to set a StartupInfo.OnAppQuit handler.")
	}
	if info.OnAppQuit != nil {
		event.AddHandler("/app/finalize", func(e *event.Event) {
			info.OnAppQuit()
		})
	}
	if err := event.AddHandler("/app/init", func(e *event.Event) {
		info.OnAppStart()
	}); err != nil {
		return err
	}
	if err := window.InitWindows(info); err != nil {
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

//	driver.Run()
