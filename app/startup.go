package app

import (
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/menu"
)

type StartupInfo struct {
	driver.StartupInfo
	AppMenu menu.AppMenuTemplate
}

type StartupFunc func(*StartupInfo) error
