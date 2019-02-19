package menu

import (
	"fmt"

	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/internal/markup"
	"github.com/yossoy/exciton/internal/object"
	"github.com/yossoy/exciton/log"
)

type menuClass struct {
	event.EventHostCore
}

func (mc *menuClass) GetTarget(id string, parent event.EventTarget) event.EventTarget {
	log.PrintDebug("MenuClass: GetTarget() : %q", id)
	itm := object.Menus.Get(object.ObjectKey(id))
	log.PrintDebug("==> %v", id)
	if itm == nil {
		return nil
	}
	return itm.(*MenuInstance)
}

var MenuClass menuClass

type Owner interface {
	event.EventTarget
	EventPath(fragments ...string) string
	EventPath2(fragments1 []string, fragments2 []string) string
}

type MenuInstance struct {
	event.EventTarget
	markup.Buildable
	builder markup.Builder
	mounted markup.RenderResult
	uuid    string
	owner   Owner
}

func (m *MenuInstance) TargetID() string {
	return m.uuid
}

func (m *MenuInstance) GetEventSlot(name string) *event.Slot {
	return nil
}

func (m *MenuInstance) Host() event.EventHost {
	return &MenuClass
}

func (m *MenuInstance) ParentTarget() event.EventTarget {
	return m.owner
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
	result := event.EmitWithResult(m, "updateDiffSetHandler", event.NewValue(ds))
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
	object.Menus.Put(uid, m)

	result := event.EmitWithResult(m, "new", event.NewValue(nil))
	if result.Error() != nil {
		object.Menus.Delete(uid)
		return nil, result.Error()
	}

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
	firstItem := true
	addSeparator := false
	for _, m := range menu {
		if m.Hidden {
			continue
		}
		if m.Separator {
			if !firstItem {
				addSeparator = true
			}
			continue
		}
		if addSeparator {
			items = append(items, html.HorizontalRule())
			addSeparator = false
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
		firstItem = false
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

func InitEvents(owner event.EventHost) {
	event.InitHost(&MenuClass, "menu", owner)
	MenuClass.AddHandler("finalize", func(e *event.Event) {
		m, ok := e.Target.(*MenuInstance)
		if !ok {
			return
		}
		object.Menus.Delete(object.ObjectKey(m.uuid))
	})
	markup.InitEvents(&MenuClass)
}
