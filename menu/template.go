package menu

import (
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"
)

type AppMenuTemplate []AppMenuItemTemplate

type AppMenuItemTemplate struct {
	Label   string
	Hidden  bool
	Role    MenuRole
	SubMenu MenuTemplate
}

type MenuTemplate []ItemTemplate

type ItemTemplate struct {
	Label      string
	Acclerator string
	Hidden     bool
	Role       MenuRole
	Handler    func(e *html.MouseEvent)
	SubMenu    MenuTemplate
	Separator  bool
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
