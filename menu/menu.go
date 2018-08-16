package menu

import (
	"fmt"

	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"
	ievent "github.com/yossoy/exciton/internal/event"
	"github.com/yossoy/exciton/internal/object"
	"github.com/yossoy/exciton/log"
	"github.com/yossoy/exciton/markup"
)

type MenuInstance struct {
	builder   markup.Builder
	mounted   markup.RenderResult
	uuid      string
	eventRoot string
}

func (m *MenuInstance) Builder() markup.Builder {
	return m.builder
}

func (m *MenuInstance) EventRoot() string {
	return m.eventRoot
}

func (m *MenuInstance) requestAnimationFrame() {
	//	go func() {
	m.builder.ProcRequestAnimationFrame()
	//}()
	log.PrintInfo("called requestAnimationFrame")
}

func (m *MenuInstance) updateDiffSetHandler(ds *markup.DiffSet) {
	result := ievent.EmitWithResult(m.eventRoot+"/menu/"+m.uuid+"/updateDiffSetHandler", event.NewValue(ds))
	if result.Error() != nil {
		panic(result.Error())
	}
	var ret bool
	if e := result.Value().Decode(&ret); e != nil {
		panic(e)
	}
	if !ret {
		panic("invalid /menu/" + m.uuid + "/updateDiffSetHandler" + " results")
	}
}

func newMenu(eventRoot string) (*MenuInstance, error) {
	uid := object.Menus.NewKey()

	result := ievent.EmitWithResult(eventRoot+"/menu/"+uid+"/new", event.NewValue(nil))
	if result.Error() != nil {
		return nil, result.Error()
	}

	m := &MenuInstance{
		uuid:      uid,
		eventRoot: eventRoot,
	}
	object.Menus.Put(uid, m)

	return m, nil
}

func newInstance(eventRoot string, component markup.RenderResult) (*MenuInstance, error) {
	m, err := newMenu(eventRoot)
	if err != nil {
		return nil, err
	}

	m.mounted = component
	m.builder = markup.NewAsyncBuilder(eventRoot+"/menu/"+m.uuid, m.requestAnimationFrame, m.updateDiffSetHandler)
	m.builder.RenderBody(component)

	return m, nil
}

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

func InitMenus(gg event.Group) error {
	err := gg.AddHandler("/menu/:id/finalize", func(e *event.Event) {
		id := e.Params["id"]
		_, _, err := object.Menus.Delete(object.ObjectKey(id))
		if err != nil {
			panic(err)
		}
	})

	return err
}
