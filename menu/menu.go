package menu

import (
	"fmt"
	"strings"

	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/internal/markup"
	"github.com/yossoy/exciton/internal/object"
	"github.com/yossoy/exciton/log"
)

type menuClass struct {
	event.EventHostCore
}

func (mc *menuClass) GetTarget(id string, parent event.EventTarget) event.EventTarget {
	itm := object.Menus.Get(object.ObjectKey(id))
	if itm == nil {
		return nil
	}
	return itm.(*MenuInstance)
}

type menuEventSlot struct {
	core event.Slot
}

func (s *menuEventSlot) Core() *event.Slot {
	return &s.core
}

func (s *menuEventSlot) Bind(h event.Handler) {
	s.core.Bind(h)
}

func (s *menuEventSlot) BindWithResult(h event.HandlerWithResult) {
	s.core.BindWithResult(h)
}

func (s *menuEventSlot) IsEnabled() bool {
	return s.core.IsEnabled()
}

func (s *menuEventSlot) SetValidateEnabledHandler(validator func(name string) bool) {
	s.core.SetValidateEnabledHandler(validator)
}

var MenuClass menuClass

type Owner interface {
	event.EventTarget
	event.EventTargetWithScopedNameResolver
}

type MenuInstance struct {
	// markup.Buildable
	// builder            markup.Builder
	mounted            markup.RenderResult
	uuid               string
	owner              Owner
	scopedNameResolver event.EventTargetWithScopedNameResolver
	items              map[string]*menuItem
	lastMenuID         int
}

func (m *MenuInstance) TargetID() string {
	return m.uuid
}

func (m *MenuInstance) Host() event.EventHost {
	return &MenuClass
}

func (m *MenuInstance) ParentTarget() event.EventTarget {
	return m.owner
}

func (m *MenuInstance) GetTargetByScopedName(scopedName string) (event.EventTarget, string) {
	if m.scopedNameResolver == nil {
		return nil, scopedName
	}
	return m.scopedNameResolver.GetTargetByScopedName(scopedName)
}

// func (m *MenuInstance) Builder() markup.Builder {
// 	return m.builder
// }

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
		panic(fmt.Sprintf("invalid results: &v", result.Value()))
	}
}

func (m *MenuInstance) cnvAppMenuTemplate(mt AppMenuTemplate) ([]*appMenuItem, error) {
	r := make([]*appMenuItem, 0, len(mt))
	for _, mm := range mt {
		if mm.Hidden {
			continue
		}
		rr := &appMenuItem{
			Label: mm.Label,
			Role:  mm.Role,
		}
		if len(mm.SubMenu) > 0 {
			sm, err := m.cnvMenuTemplate(mm.SubMenu)
			if err != nil {
				return nil, err
			}
			rr.SubMenu = sm
		}
		r = append(r, rr)
	}
	return r, nil
}

func (m *MenuInstance) cnvMenuTemplate(mt MenuTemplate) ([]*menuItem, error) {
	r := make([]*menuItem, 0, len(mt))
	for _, mm := range mt {
		if mm.Hidden {
			continue
		}
		if strings.HasPrefix(mm.ID, "exciton_") {
			return nil, fmt.Errorf("menuitem: ID cannot start 'exciton'")
		}
		rr := &menuItem{
			ID:         mm.ID,
			Label:      mm.Label,
			Acclerator: mm.Acclerator,
			Role:       mm.Role,
			Separator:  mm.Separator,
			Action:     mm.Action,
			handler:    mm.Handler,
		}
		if rr.ID == "" {
			m.lastMenuID++
			rr.ID = fmt.Sprintf("exciton_id_%d", m.lastMenuID)
		}
		if len(mm.SubMenu) > 0 {
			sm, err := m.cnvMenuTemplate(mm.SubMenu)
			if err != nil {
				return nil, err
			}
			rr.SubMenu = sm
		}
		if mm.Action != "" && mm.Handler != nil {
			return nil, fmt.Errorf("menuitem: cannot specify both Action and Handler")
		}
		if mm.ID != "" || mm.Action != "" || mm.Handler != nil {
			if m.items == nil {
				m.items = make(map[string]*menuItem)
			}
			m.items[rr.ID] = rr
		}
		r = append(r, rr)
	}
	return r, nil
}

func newMenu(owner Owner, scopedNameResolver event.EventTargetWithScopedNameResolver) (*MenuInstance, error) {
	uid := object.Menus.NewKey()

	m := &MenuInstance{
		uuid:               uid,
		owner:              owner,
		scopedNameResolver: scopedNameResolver,
	}
	object.Menus.Put(uid, m)

	result := event.EmitWithResult(m, "new", event.NewValue(nil))
	if result.Error() != nil {
		object.Menus.Delete(uid)
		return nil, result.Error()
	}

	return m, nil
}

func newPopupMenu(owner Owner, scopedNameResolver event.EventTargetWithScopedNameResolver, templ MenuTemplate) (*MenuInstance, error) {
	uid := object.Menus.NewKey()

	m := &MenuInstance{
		uuid:               uid,
		owner:              owner,
		scopedNameResolver: scopedNameResolver,
	}
	mm, err := m.cnvMenuTemplate(templ)
	if err != nil {
		return nil, err
	}
	object.Menus.Put(uid, m)
	result := event.EmitWithResult(m, "newPopupMenu", event.NewValue(mm))
	if result.Error() != nil {
		object.Menus.Delete(uid)
		return nil, result.Error()
	}

	return m, nil
}

func newAppMenu(owner Owner, scopedNameResolver event.EventTargetWithScopedNameResolver, templ AppMenuTemplate) (*MenuInstance, error) {
	uid := object.Menus.NewKey()

	m := &MenuInstance{
		uuid:               uid,
		owner:              owner,
		scopedNameResolver: scopedNameResolver,
	}
	am, err := m.cnvAppMenuTemplate(templ)
	if err != nil {
		return nil, err
	}
	object.Menus.Put(uid, m)

	result := event.EmitWithResult(m, "newApplicationMenu", event.NewValue(am))
	if result.Error() != nil {
		object.Menus.Delete(uid)
		return nil, result.Error()
	}

	return m, nil
}

// func newInstance(owner Owner, scopedNameResolver event.EventTargetWithScopedNameResolver, component markup.RenderResult) (*MenuInstance, error) {
// 	m, err := newMenu(owner, scopedNameResolver)
// 	if err != nil {
// 		return nil, err
// 	}

// 	m.mounted = component
// 	m.builder = markup.NewAsyncBuilder(m)
// 	m.builder.RenderBody(component)

// 	return m, nil
// }

// func toPopupMenuSub(menu MenuTemplate) ([]markup.MarkupOrChild, error) {
// 	var items []markup.MarkupOrChild
// 	firstItem := true
// 	addSeparator := false
// 	for _, m := range menu {
// 		if m.Hidden {
// 			continue
// 		}
// 		if m.Separator {
// 			if !firstItem {
// 				addSeparator = true
// 			}
// 			continue
// 		}
// 		if addSeparator {
// 			items = append(items, html.HorizontalRule())
// 			addSeparator = false
// 		}
// 		if m.Label == "" && m.Role == "" {
// 			return nil, fmt.Errorf("menu need Label or Role")
// 		}
// 		var mitems []markup.MarkupOrChild
// 		if m.Label != "" {
// 			mitems = append(mitems, markup.AttrApplyer{Name: "label", Value: m.Label})
// 		}
// 		if m.Role != "" {
// 			mitems = append(mitems, markup.DataApplyer{Name: "menuRole", Value: string(m.Role)})
// 		}
// 		if m.Acclerator != "" {
// 			mitems = append(mitems, markup.DataApplyer{Name: "menuAcclerator", Value: m.Acclerator})
// 		}
// 		if m.Handler != nil {
// 			//TODO: modify event type
// 			mitems = append(mitems, html.OnClick(m.Handler))
// 		}
// 		if len(m.SubMenu) > 0 {
// 			smitems, err := toPopupMenuSub(m.SubMenu)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if len(smitems) > 0 {
// 				mitems = append(mitems, smitems...)
// 			}
// 		}
// 		if len(m.SubMenu) > 0 || roleIsMenuedRole(m.Role) {
// 			items = append(items, markup.MustTag("menu", mitems))
// 		} else {
// 			items = append(items, markup.MustTag("menuitem", mitems))
// 		}
// 		firstItem = false
// 	}
// 	return items, nil
// }

// func toAppMenu(menu AppMenuTemplate) (markup.RenderResult, error) {
// 	var items []markup.MarkupOrChild
// 	for _, m := range menu {
// 		if m.Hidden {
// 			continue
// 		}
// 		var mitems []markup.MarkupOrChild
// 		if m.Label == "" {
// 			return nil, fmt.Errorf("Application Menu need Label")
// 		}
// 		mitems = append(mitems, markup.AttrApplyer{Name: "label", Value: m.Label})
// 		if m.Role != "" {
// 			mitems = append(mitems, markup.DataApplyer{Name: "menuRole", Value: string(m.Role)})
// 		}
// 		if m.SubMenu != nil {
// 			smitems, err := toPopupMenuSub(m.SubMenu)
// 			if err != nil {
// 				return nil, err
// 			}
// 			mitems = append(mitems, smitems...)
// 		}
// 		if m.SubMenu != nil || roleIsMenuedRole(m.Role) {
// 			items = append(items, markup.MustTag("menu", mitems))
// 		} else {
// 			items = append(items, markup.MustTag("menuitem", mitems))
// 		}
// 	}
// 	return markup.Tag("menu", items)
// }

func getMenuState(e *event.Event, callback event.ResponceCallback) {
	mi, ok := e.Target.(*MenuInstance)
	if !ok {
		panic(false)
	}
	var id string
	err := e.Argument.Decode(&id)
	if err != nil {
		callback(event.NewErrorResult(err))
		return
	}
	itm, ok := mi.items[id]
	if !ok {
		callback(event.NewErrorResult(fmt.Errorf("menuitem: %q not found", id)))
		return
	}
	if itm.Action != "" {
		et, name := mi.scopedNameResolver.GetTargetByScopedName(itm.Action)
		if et == nil {
			callback(event.NewErrorResult(fmt.Errorf("menu: Action %q not found", itm.Action)))
			return
		}
		log.PrintDebug("name = %q", name)

	}
	callback(event.NewValueResult(event.NewValue(true)))
}

func emitMenu(e *event.Event) error {
	mi, ok := e.Target.(*MenuInstance)
	if !ok {
		panic(false)
	}
	var id string
	err := e.Argument.Decode(&id)
	if err != nil {
		return err
	}
	itm, ok := mi.items[id]
	if !ok {
		return fmt.Errorf("menuitem: %q not found", id)
	}
	if itm.Action != "" {
		et, name := mi.scopedNameResolver.GetTargetByScopedName(itm.Action)
		if et == nil {
			err = fmt.Errorf("menu: Action %q not found", itm.Action)
		} else {
			err = event.Emit(et, name, e.Argument)
		}
		log.PrintError("err = %v", err)
	} else if itm.handler != nil {
		err = itm.handler(e)
	}
	if err != nil {
		return err
	}
	return nil
}
func isEnableMenu(e *event.Event, callback event.ResponceCallback) {
	mi, ok := e.Target.(*MenuInstance)
	if !ok {
		panic(false)
	}
	var id string
	err := e.Argument.Decode(&id)
	if err != nil {
		callback(event.NewErrorResult(err))
		return
	}
	itm, ok := mi.items[id]
	if !ok {
		callback(event.NewErrorResult(fmt.Errorf("menuitem: %q not found", id)))
		return
	}
	if itm.Action != "" {
		et, name := mi.scopedNameResolver.GetTargetByScopedName(itm.Action)
		if et == nil {
			callback(event.NewErrorResult(fmt.Errorf("menu: Action %q not found", itm.Action)))
		} else {
			enabled, err := event.IsEnableEvent(et, name)
			if err != nil {
				callback(event.NewErrorResult(err))
			} else {
				callback(event.NewValueResult(event.NewValue(enabled)))
			}
		}
	} else if itm.handler != nil {
		callback(event.NewValueResult(event.NewValue(true)))
	}
}

func InitEvents(owner event.EventHost) {
	event.InitHost(&MenuClass, "menu", owner)
	MenuClass.AddHandler("finalize", func(e *event.Event) error {
		m, ok := e.Target.(*MenuInstance)
		if !ok {
			return fmt.Errorf("invalid target: %v", e.Target)
		}
		object.Menus.Delete(object.ObjectKey(m.uuid))
		return nil
	})
	MenuClass.AddHandler("emit", emitMenu)
	MenuClass.AddHandlerWithResult("is-item-enabled", isEnableMenu)
	MenuClass.AddHandlerWithResult("status", func(e *event.Event, callback event.ResponceCallback) {
		var id string
		e.Argument.Decode(&id)
		callback(event.NewErrorResult(fmt.Errorf("menu: not implement yet")))
	})
	markup.InitEvents(&MenuClass)
}
