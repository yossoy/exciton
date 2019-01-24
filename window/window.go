package window

import (
	"errors"
	"fmt"

	"github.com/yossoy/exciton/html"

	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/geom"
	ievent "github.com/yossoy/exciton/internal/event"
	"github.com/yossoy/exciton/internal/object"
	"github.com/yossoy/exciton/log"
	"github.com/yossoy/exciton/markup"
)

// UserData is window binded data
type UserData interface{}

type Owner interface {
	ID() string
	PreferredLanguages() []string
	EventPath(fragments ...string) string
	EventPath2(fragments1 []string, fragments2 []string) string
	URLBase() string
}

// Window is browser window
type Window struct {
	markup.Buildable
	owner             Owner
	ID                string
	UserData          UserData
	builder           markup.Builder
	isReady           bool
	OnClosed          func(e *event.Event)
	OnResize          func(width, height float64)
	OnKeyDown         func(w *Window, e *html.KeyboardEvent)
	OnKeyUp           func(w *Window, e *html.KeyboardEvent)
	mountRenderResult markup.RenderResult
	title             string
	lang              string
	cachedHTML        []byte
}

var activeWindow *Window

func ActiveWindow() *Window {
	return activeWindow
}

func (w *Window) Owner() Owner {
	return w.owner
}

func (w *Window) Builder() markup.Builder {
	return w.builder
}

func (w *Window) EventPath(fragments ...string) string {
	return w.owner.EventPath2([]string{"window", w.ID}, fragments)
}

func (w *Window) EventPath2(fragments1 []string, fragments2 []string) string {
	f := make([]string, 2, 2+len(fragments1))
	f[0] = "window"
	f[1] = w.ID
	f = append(f, fragments1...)
	return w.owner.EventPath2(f, fragments2)
}

func (w *Window) RequestAnimationFrame() {
	ievent.Emit(w.EventPath("requestAnimationFrame"), nil)
}

func (w *Window) UpdateDiffSetHandler(ds *markup.DiffSet) {
	ievent.Emit(w.EventPath("updateDiffSetHandler"), event.NewValue(ds))
}

func (w *Window) onReady() {
	log.PrintInfo("window ready: %v\n", w)
	if w.isReady {
		return
	}
	w.isReady = true
	if w.mountRenderResult != nil {
		w.builder.RenderBody(w.mountRenderResult)
	}
}

func (w *Window) onRequestAnimationFrame(tick float64) {
	log.PrintInfo("onRequestAnimationFrame(%f", tick)
	w.builder.ProcRequestAnimationFrame()
}

func (w *Window) Mount(c markup.RenderResult) error {
	if c == nil {
		return errors.New("Windows is already mount")
	}
	if w.mountRenderResult != nil {
		return errors.New("Windows is already mount")
	}
	w.mountRenderResult = c
	if w.isReady {
		w.builder.RenderBody(w.mountRenderResult)
	}
	return nil
}

// WindowConfig is a struct that describes a window.
type WindowConfig struct {
	ID              string              `json:"id"`
	Title           string              `json:"title,omitempty"`
	Position        *geom.Point         `json:"position,omitempty"`
	Size            *geom.Size          `json:"size,omitempty"`
	MinSize         *geom.Size          `json:"minSize,omitempty"`
	MaxSize         *geom.Size          `json:"maxSize,omitempty"`
	BackgroundColor string              `json:"backgroundColor,omitempty"`
	FixedSize       bool                `json:"fixedSize"`
	NoClosable      bool                `json:"noClosable"`
	NoMinimizable   bool                `json:"noMinimizable"`
	Lang            string              `json:"lang"`
	URL             string              `json:"url"`
	OnCreate        func(*Window) error `json:"-"`
}

const (
	stdWinWidth  float64 = 640.0
	stdWinHeight float64 = 480.0
)

func getWindowByID(id string) *Window {
	itm := object.Windows.Get(object.ObjectKey(id))
	if itm == nil {
		return nil
	}
	return itm.(*Window)
}

func onWindowClosed(e *event.Event) {
	w := getWindowByID(e.Params["id"])
	if w == nil {
		return
	}
	if activeWindow == w {
		activeWindow = nil
	}
	if w.OnClosed != nil {
		w.OnClosed(e)
	}
}

func onWindowResized(e *event.Event) {
	w := getWindowByID(e.Params["id"])
	if w == nil {
		return
	}
	sz := geom.Size{}
	if err := e.Argument.Decode(&sz); err != nil {
		log.PrintError("/window/:id/resize: parameter decode failed: %#v", e.Argument)
		return
	}
	log.PrintDebug("Window: resized (%#v)", sz)
	if w.OnResize != nil {
		w.OnResize(sz.Width, sz.Height)
	}
}

func onWindowKeyDown(e *event.Event) {
	w := getWindowByID(e.Params["id"])
	if w == nil {
		return
	}
	var ke html.KeyboardEvent
	err := e.Argument.Decode(&ke)
	if err != nil {
		log.PrintError("/window/:id/keydown: parameter decode failed: %#v", e.Argument)
		return
	}
	log.PrintDebug("Window: keydown (%#v)", ke)
	if w.OnKeyDown != nil {
		w.OnKeyDown(w, &ke)
	}
}

func onWindowKeyUp(e *event.Event) {
	w := getWindowByID(e.Params["id"])
	if w == nil {
		return
	}
	var ke html.KeyboardEvent
	err := e.Argument.Decode(&ke)
	if err != nil {
		log.PrintError("/window/:id/keyup: parameter decode failed: %#v", e.Argument)
		return
	}
	log.PrintDebug("Window: keyup (%#v)", ke)
	if w.OnKeyUp != nil {
		w.OnKeyUp(w, &ke)
	}
}

func onWindowFocus(e *event.Event) {
	w := getWindowByID(e.Params["id"])
	if w == nil {
		return
	}
	activeWindow = w
}

func onWindowBlur(e *event.Event) {
	w := getWindowByID(e.Params["id"])
	if w == nil {
		return
	}
	if activeWindow == w {
		activeWindow = nil
	}
}

type changeRouteArg struct {
	Route string `json:"route"`
}

func onChangeRoute(e *event.Event) {
	w := getWindowByID(e.Params["id"])
	var arg changeRouteArg
	if err := e.Argument.Decode(&arg); err != nil {
		panic(err)
	}
	log.PrintDebug("onChangeRoute: %q", arg.Route)
	w.Builder().OnRedirect(arg.Route)
}

func InitWindows(si *driver.StartupInfo) error {
	appg := si.AppEventRoot
	if err := initHTML(si); err != nil {
		return err
	}
	if err := appg.AddHandler("/window/:id/finalize", func(e *event.Event) {
		id := e.Params["id"]
		log.PrintInfo("finalized: %q", id)
		wobj, empty, err := object.Windows.Delete(object.ObjectKey(id))
		if err != nil {
			log.PrintInfo("finalize: %q", err)
			return
		}
		win := wobj.(*Window)
		ievent.Emit("/app/finalizedWindow", event.NewValue(win))
		if empty {
			ievent.Emit("/app/window-all-closed", nil)
		}
	}); err != nil {
		return err
	}
	if err := appg.AddHandler("/window/:id/changeRoute", onChangeRoute); err != nil {
		return err
	}
	if err := appg.AddHandler("/window/:id/closed", onWindowClosed); err != nil {
		return err
	}
	if err := appg.AddHandler("/window/:id/resize", onWindowResized); err != nil {
		return err
	}
	if err := appg.AddHandler("/window/:id/keydown", onWindowKeyDown); err != nil {
		return err
	}
	if err := appg.AddHandler("/window/:id/keyup", onWindowKeyUp); err != nil {
		return err
	}
	if err := appg.AddHandler("/window/:id/focus", onWindowFocus); err != nil {
		return err
	}
	if err := ievent.AddHandler("/window/:id/blur", onWindowBlur); err != nil {
		return err
	}

	if err := appg.AddHandler("/window/:id/onRequestAnimationFrame", func(e *event.Event) {
		w := getWindowByID(e.Params["id"])
		var tick float64
		if err := e.Argument.Decode(&tick); err != nil {
			panic(err)
		}
		w.onRequestAnimationFrame(tick)
	}); err != nil {
		return err
	}

	if err := appg.AddHandler("/window/:id/ready", func(e *event.Event) {
		w := getWindowByID(e.Params["id"])
		w.onReady()
	}); err != nil {
		return err
	}
	log.PrintInfo("initß ok\n")
	return nil
}

// NewWindow create new browser window
func NewWindow(owner Owner, cfg *WindowConfig) (*Window, error) {
	if cfg.Size == nil {
		cfg.Size = &geom.Size{Width: stdWinWidth, Height: stdWinHeight}
	}
	id := object.Windows.NewKey()
	cfg.ID = id
	if cfg.URL == "" {
		appURLBase := owner.URLBase()
		cfg.URL = fmt.Sprintf("%s"+appURLBase+"/window/%s/", driver.BaseURL, id)
	}

	w := &Window{
		owner: owner,
		ID:    id,
		title: cfg.Title,
		lang:  cfg.Lang,
	}
	object.Windows.Put(id, w)
	r := ievent.EmitWithResult(w.EventPath("new"), event.NewValue(cfg))
	if r.Error() != nil {
		object.Windows.Delete(id)
		return nil, r.Error()
	}
	w.builder = markup.NewAsyncBuilder(w)
	if cfg.OnCreate != nil {
		cfg.OnCreate(w)
	}

	return w, nil
}

func GetWindowFromEventTarget(e *markup.EventTarget) (*Window, error) {
	if e.WindowID == "" {
		return nil, errors.New("invalid event target")
	}
	w := getWindowByID(e.WindowID)
	if w == nil {
		return nil, errors.New("window not found")
	}
	return w, nil
}
