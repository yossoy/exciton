package app

import (
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
)

type StartupInfo struct {
	driver.StartupInfo
	AppMenu     menu.AppMenuTemplate
	OnAppStart  func(*App, *StartupInfo) error
	OnAppQuit   func()
	OnNewWindow func(*App, *window.WindowConfig) (markup.RenderResult, error)
}

type StartupFunc func(*StartupInfo) error
