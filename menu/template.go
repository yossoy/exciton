package menu

import (
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"
)

type AppMenuTemplate []AppMenuItemTemplate

type AppMenuItemTemplate struct {
	Label   string       `json:"label"`
	Hidden  bool         `json:"hidden"`
	Role    MenuRole     `json:"role"`
	SubMenu MenuTemplate `json:"subMenu"`
}

type MenuTemplate []ItemTemplate

type ItemTemplate struct {
	Label           string       `json:"label"`
	Acclerator      string       `json:"acclerator"`
	Hidden          bool         `json:"hidden"`
	Role            MenuRole     `json:"role"`
	EventSignalName string       `json:"eventSignalName"`
	SubMenu         MenuTemplate `json:"subMenu"`
	Separator       bool         `json:"separator"`
	Handler         func(e *html.MouseEvent) error
}

type MenuEvent struct {
	View   event.EventTarget
	Target event.EventTarget
}

const (
	AppMenuLabel = "*Appname*"
)

var (
	Separator = ItemTemplate{Separator: true}
)
