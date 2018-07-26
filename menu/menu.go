package menu

import (
	"fmt"

	"github.com/yossoy/exciton/geom"
	"github.com/yossoy/exciton/window"

	"github.com/yossoy/exciton/html"
	imenu "github.com/yossoy/exciton/internal/menu"
	"github.com/yossoy/exciton/markup"
)

//TODO: section support for gtk

func toPopupMenuSub(menu MenuTemplate) ([]markup.MarkupOrChild, error) {
	var items []markup.MarkupOrChild
	for sidx, s := range menu {
		firstItem := true
		for _, m := range s {
			if m.Hidden {
				continue
			}
			if m.Label == "" && m.Role == "" {
				return nil, fmt.Errorf("menu need Label or Role")
			}
			var mitems []markup.MarkupOrChild
			if m.Label != "" {
				mitems = append(mitems, markup.Attribute("label", m.Label))
			}
			if m.Role != "" {
				mitems = append(mitems, markup.Data("menuRole", string(m.Role)))
			}
			if m.Acclerator != "" {
				mitems = append(mitems, markup.Data("menuAcclerator", m.Acclerator))
			}
			if m.Handler != nil {
				//TODO: modify event type
				mitems = append(mitems, html.OnClick(m.Handler))
			}
			if firstItem && sidx != 0 {
				items = append(items, html.HorizontalRule())
				firstItem = false
			}
			if len(m.SubMenu) > 0 {
				smitems, err := toPopupMenuSub(m.SubMenu)
				if err != nil {
					return nil, err
				}
				if len(smitems) > 0 {
					mitems = append(mitems, smitems...)
				}
			}
			if len(m.SubMenu) > 0 || roleIsMenuedRole(m.Role) {
				items = append(items, markup.Tag("menu", mitems...))
			} else {
				items = append(items, markup.Tag("menuitem", mitems...))
			}
		}
	}
	return items, nil
}

func toAppMenu(menu AppMenuTemplate) (markup.RenderResult, error) {
	var items []markup.MarkupOrChild
	for _, m := range menu {
		if m.Hidden {
			continue
		}
		var mitems []markup.MarkupOrChild
		if m.Label == "" {
			return nil, fmt.Errorf("Application Menu need Label")
		}
		mitems = append(mitems, markup.Attribute("label", m.Label))
		if m.Role != "" {
			mitems = append(mitems, markup.Data("menuRole", string(m.Role)))
		}
		if m.SubMenu != nil {
			smitems, err := toPopupMenuSub(m.SubMenu)
			if err != nil {
				return nil, err
			}
			mitems = append(mitems, smitems...)
		}
		if m.SubMenu != nil || roleIsMenuedRole(m.Role) {
			items = append(items, markup.Tag("menu", mitems...))
		} else {
			items = append(items, markup.Tag("menuitem", mitems...))
		}
	}
	return markup.Tag("menu", items...), nil
}

func SetApplicationMenu(menu AppMenuTemplate) error {
	r, err := toAppMenu(menu)
	if err != nil {
		return err
	}
	mi, err := imenu.New(r)
	if err != nil {
		return err
	}
	return imenu.SetApplicationMenu(mi)
}

func PopupMenu(menu MenuTemplate, mousePt geom.Point, w *window.Window) error {
	m, err := toPopupMenuSub(menu)
	if err != nil {
		return err
	}
	r := markup.Tag("menu", m...)
	mi, err := imenu.New(r)
	if err != nil {
		return err
	}
	return mi.Popup(mousePt, w)
}
