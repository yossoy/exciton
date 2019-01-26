package menu

import (
	"fmt"

	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"
	ievent "github.com/yossoy/exciton/internal/event"
	"github.com/yossoy/exciton/internal/markup"
	"github.com/yossoy/exciton/internal/object"
	"github.com/yossoy/exciton/log"
)

type Owner interface {
	EventPath(fragments ...string) string
	EventPath2(fragments1 []string, fragments2 []string) string
}

type MenuInstance struct {
	markup.Buildable
	builder markup.Builder
	mounted markup.RenderResult
	uuid    string
	owner   Owner
}

func (m *MenuInstance) Builder() markup.Builder {
	return m.builder
}

func (m *MenuInstance) EventPath(fragments ...string) string {
	return m.owner.EventPath2([]string{"menu", m.uuid}, fragments)
}

func (m *MenuInstance) RequestAnimationFrame() {
	m.builder.ProcRequestAnimationFrame()
	log.PrintInfo("called requestAnimationFrame")
}

func (m *MenuInstance) UpdateDiffSetHandler(ds *markup.DiffSet) {
	result := ievent.EmitWithResult(m.EventPath("updateDiffSetHandler"), event.NewValue(ds))
	if result.Error() != nil {
		panic(result.Error())
	}
	var ret bool
	if e := result.Value().Decode(&ret); e != nil {
		panic(e)
	}
	if !ret {
		panic(fmt.Sprintf("invalid %q results", m.EventPath("updateDiffSetHandler")))
	}
}

func newMenu(owner Owner) (*MenuInstance, error) {
	uid := object.Menus.NewKey()

	m := &MenuInstance{
		uuid:  uid,
		owner: owner,
	}

	result := ievent.EmitWithResult(m.EventPath("new"), event.NewValue(nil))
	if result.Error() != nil {
		return nil, result.Error()
	}

	object.Menus.Put(uid, m)

	return m, nil
}

func newInstance(owner Owner, component markup.RenderResult) (*MenuInstance, error) {
	m, err := newMenu(owner)
	if err != nil {
		return nil, err
	}

	m.mounted = component
	m.builder = markup.NewAsyncBuilder(m)
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
				mitems = append(mitems, markup.AttrApplyer{Name: "label", Value: m.Label})
			}
			if m.Role != "" {
				mitems = append(mitems, markup.DataApplyer{Name: "menuRole", Value: string(m.Role)})
			}
			if m.Acclerator != "" {
				mitems = append(mitems, markup.DataApplyer{Name: "menuAcclerator", Value: m.Acclerator})
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
				items = append(items, markup.MustTag("menu", mitems))
			} else {
				items = append(items, markup.MustTag("menuitem", mitems))
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
		mitems = append(mitems, markup.AttrApplyer{Name: "label", Value: m.Label})
		if m.Role != "" {
			mitems = append(mitems, markup.DataApplyer{Name: "menuRole", Value: string(m.Role)})
		}
		if m.SubMenu != nil {
			smitems, err := toPopupMenuSub(m.SubMenu)
			if err != nil {
				return nil, err
			}
			mitems = append(mitems, smitems...)
		}
		if m.SubMenu != nil || roleIsMenuedRole(m.Role) {
			items = append(items, markup.MustTag("menu", mitems))
		} else {
			items = append(items, markup.MustTag("menuitem", mitems))
		}
	}
	return markup.Tag("menu", items)
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
