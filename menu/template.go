package menu

import (
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/markup"
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
	View   *markup.EventTarget
	Target *markup.EventTarget
}

const (
	AppMenuLabel = "*Appname*"
)

var (
	Separator = ItemTemplate{Separator: true}
)
