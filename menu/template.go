package menu

import (
	"github.com/yossoy/exciton/event"
)

type AppMenuTemplate []AppMenuItemTemplate

type AppMenuItemTemplate struct {
	Label   string
	Hidden  bool
	Role    MenuRole
	SubMenu MenuTemplate
}

type appMenuItem struct {
	Label   string      `json:"label"`
	Role    MenuRole    `json:"role"`
	SubMenu []*menuItem `json:"subMenu"`
}

type MenuTemplate []ItemTemplate

type ItemTemplate struct {
	ID         string
	Label      string
	Acclerator string
	Hidden     bool
	Role       MenuRole
	SubMenu    MenuTemplate
	Separator  bool
	Action     string
	Handler    func(e *event.Event) error
}

type menuItem struct {
	ID         string                     `json:"id"`
	Label      string                     `json:"label,omitempty"`
	Acclerator string                     `json:"acclerator,omitempty"`
	Role       MenuRole                   `json:"role,omitempty"`
	SubMenu    []*menuItem                `json:"subMenu,omitempty"`
	Separator  bool                       `json:"separator"`
	Action     string                     `json:"action,omitempty"`
	handler    func(e *event.Event) error `json:"-"`
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
