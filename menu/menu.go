package menu

import (
	"strings"

	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/internal/object"
	"github.com/yossoy/exciton/log"
	"github.com/yossoy/exciton/markup"
)

type ItemRole string

const (
	RoleAbout              ItemRole = "about"
	RoleHide                        = "hide"
	RoleHideOthers                  = "hideothers"
	RoleUnhide                      = "unhide"
	RoleFront                       = "front"
	RoleUndo                        = "undo"
	RoleRedo                        = "redo"
	RoleCut                         = "cut"
	RoleCopy                        = "copy"
	RolePaste                       = "paste"
	RoleDelete                      = "delete"
	RolePasteAndMatchStyle          = "pasteandmatchstyle"
	RoleSelectAll                   = "selectall"
	RoleStartSpeaking               = "startspeaking"
	RoleStopSpeaking                = "stopspeaking"
	RoleMinimize                    = "minimize"
	RoleClose                       = "close"
	RoleZoom                        = "zoom"
	RoleQuit                        = "quit"
	RoleToggleFullscreen            = "togglefullscreen"
)

type MenuRole string

const (
	RoleServices MenuRole = "services"
	RoleWindow            = "window"
	RoleHelp              = "help"
)

type MenuInstance struct {
	builder *markup.Builder
	mounted *markup.RenderResult
	uuid    string
}

func (m *MenuInstance) Builder() *markup.Builder {
	return m.builder
}

func (m *MenuInstance) requestAnimationFrame() {
	//	go func() {
	m.builder.ProcRequestAnimationFrame()
	//}()
	log.PrintInfo("called requestAnimationFrame")
}

func (m *MenuInstance) updateDiffSetHandler(ds *markup.DiffSet) {
	result := event.EmitWithResult("/menu/"+m.uuid+"/updateDiffSetHandler", event.NewValue(ds))
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

func SubMenu(label string, m ...markup.MarkupOrChild) *markup.RenderResult {
	if len(m) == 0 {
		//TODO: return emtpy menu?
		return nil
	}
	args := make([]markup.MarkupOrChild, 0, 1+len(m))
	args = append(args, markup.Attribute("label", label))
	args = append(args, m...)
	return markup.Tag("menu", args...)
}

func Item(label string, handler func(e *html.MouseEvent)) *markup.RenderResult {
	ll := strings.Split(label, ";")
	if len(ll) >= 2 {
		return markup.Tag("menuitem",
			markup.Attribute("label", strings.TrimSpace(ll[0])),
			markup.Data("menuAcclerator", strings.TrimSpace(ll[1])),
			html.OnClick(handler),
		)
	}
	return markup.Tag("menuitem",
		markup.Attribute("label", label),
		html.OnClick(handler),
	)
}

func RoledItem(role ItemRole) *markup.RenderResult {
	return markup.Tag("menuitem",
		markup.Data("menuRole", string(role)),
	)
}

func RoledMenu(role MenuRole, label string, childitem ...markup.MarkupOrChild) *markup.RenderResult {
	if len(childitem) == 0 {
		return markup.Tag("menu",
			markup.Attribute("label", label),
			markup.Data("menuRole", string(role)),
		)
	}
	args := make([]markup.MarkupOrChild, 0, 2+len(childitem))
	args = append(args,
		markup.Attribute("label", label),
		markup.Data("menuRole", string(role)),
	)
	args = append(args, childitem...)
	return markup.Tag("menu", args...)
}

func Separator() *markup.RenderResult {
	return html.HorizontalRule()
}

func newMenu() (*MenuInstance, error) {
	uid := object.Menus.NewKey()

	result := event.EmitWithResult("/menu/"+uid+"/new", event.NewValue(nil))
	if result.Error() != nil {
		return nil, result.Error()
	}

	m := &MenuInstance{
		uuid: uid,
	}
	object.Menus.Put(uid, m)

	return m, nil
}

func New(component *markup.RenderResult) (*MenuInstance, error) {
	m, err := newMenu()
	if err != nil {
		return nil, err
	}

	m.mounted = component
	m.builder = markup.NewAsyncBuilder("/menu/"+m.uuid, m.requestAnimationFrame, m.updateDiffSetHandler)
	m.builder.RenderBody(component)

	return m, nil
}

func MustNew(component *markup.RenderResult) *MenuInstance {
	m, err := New(component)
	if err != nil {
		panic(err)
	}
	return m
}

func InitMenus() error {
	err := event.AddHandler("/menu/:id/finalize", func(e *event.Event) {
		id := e.Params["id"]
		_, _, err := object.Menus.Delete(object.ObjectKey(id))
		if err != nil {
			panic(err)
		}
	})
	return err
}
