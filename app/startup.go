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
	OnAppStart  func(*StartupInfo) error
	OnAppQuit   func()
	OnNewWindow func(cfg *window.WindowConfig) (markup.RenderResult, error)
}

type StartupFunc func(*StartupInfo) error
