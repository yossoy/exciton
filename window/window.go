package window

import (
	"errors"
	"fmt"

	"github.com/yossoy/exciton/lang"

	"github.com/yossoy/exciton/html"

	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/geom"
	"github.com/yossoy/exciton/internal/markup"
	"github.com/yossoy/exciton/internal/object"
	"github.com/yossoy/exciton/log"
)

type windowClass struct {
	event.EventHostCore
}

func (wc *windowClass) GetTarget(id string, parent event.EventTarget) event.EventTarget {
	w := getWindowByID(id)
	if w == nil {
		return nil
	}
	return w
}

var WindowClass windowClass

func onWindowFinalize(e *event.Event) {
	w, ok := e.Target.(*Window)
	if !ok {
		panic(fmt.Sprintf("invalid target: %v", e.Target))
		return
	}
	_, empty, err := object.Windows.Delete(object.ObjectKey(w.ID))
	if err != nil {
		log.PrintError("invalid window: %v", w)
		return
	}
	// TODO: emit to AppClass?
	event.Emit(w.Parent(), "finalizeWindow", event.NewValue(w))
	if empty {
		event.Emit(w.Parent(), "window-all-closed", nil)
	}
}

func InitEvents(owner event.EventHost) {
	event.InitHost(&WindowClass, "window", owner)
	WindowClass.AddHandler("finalize", onWindowFinalize)
	WindowClass.AddHandler("changeRoute", onChangeRoute)
	WindowClass.AddHandler("closed", onWindowClosed)
	WindowClass.AddHandler("resize", onWindowResized)
	WindowClass.AddHandler("keydown", onWindowKeyDown)
	WindowClass.AddHandler("keyup", onWindowKeyUp)
	WindowClass.AddHandler("focus", onWindowFocus)
	WindowClass.AddHandler("blur", onWindowBlur)
	WindowClass.AddHandler("onRequestAnimationFrame", func(e *event.Event) {
		w, ok := e.Target.(*Window)
		if !ok {
			panic(fmt.Sprintf("invalid target: %v", e.Target))
		}
		var tick float64
		if err := e.Argument.Decode(&tick); err != nil {
			panic(err)
		}
		w.onRequestAnimationFrame(tick)
	})
	WindowClass.AddHandler("ready", func(e *event.Event) {
		w, ok := e.Target.(*Window)
		if !ok {
			panic(fmt.Sprintf("invalid target: %v", e.Target))
		}
		w.onReady() // TODO
	})
	markup.InitEvents(&WindowClass)

	// TODO: child event
}

// UserData is window binded data
type UserData interface{}

type Owner interface {
	event.EventTarget
	//	ID() string
	PreferredLanguages() lang.PreferredLanguages
	URLBase() string
}

// Window is browser window
type Window struct {
	//	markup.Buildable
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

func (w *Window) GetEventSignal(name string) *event.Signal {
	// TODO: add event slot?
	return nil
}

func (w *Window) Parent() event.EventTarget {
	return w.Owner()
}

func (w *Window) Host() event.EventHost {
	return &WindowClass
}

func (w *Window) TargetID() string {
	return w.ID
}

func (w *Window) ParentTarget() event.EventTarget {
	return w.owner
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

func (w *Window) RequestAnimationFrame() {
	event.Emit(w, "requestAnimationFrame", nil)
}

func (w *Window) UpdateDiffSetHandler(ds *markup.DiffSet) {
	event.Emit(w, "updateDiffSetHandler", event.NewValue(ds))
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
	w, ok := e.Target.(*Window)
	if !ok {
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
	w, ok := e.Target.(*Window)
	if !ok {
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
	w, ok := e.Target.(*Window)
	if !ok {
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
	w, ok := e.Target.(*Window)
	if !ok {
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
	w, ok := e.Target.(*Window)
	if !ok {
		return
	}
	activeWindow = w
}

func onWindowBlur(e *event.Event) {
	w, ok := e.Target.(*Window)
	if !ok {
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
	w, ok := e.Target.(*Window)
	if !ok {
		return
	}
	var arg changeRouteArg
	if err := e.Argument.Decode(&arg); err != nil {
		panic(err)
	}
	log.PrintDebug("onChangeRoute: %q", arg.Route)
	w.Builder().OnRedirect(arg.Route)
}

func InitWindows(si *driver.StartupInfo) error {
	// appg := si.AppEventRoot
	if err := initHTML(si); err != nil {
		return err
	}
	// if err := appg.AddHandler("/window/:id/finalize", func(e *event.Event) {
	// 	id := e.Params["id"]
	// 	log.PrintInfo("finalized: %q", id)
	// 	wobj, empty, err := object.Windows.Delete(object.ObjectKey(id))
	// 	if err != nil {
	// 		log.PrintInfo("finalize: %q", err)
	// 		return
	// 	}
	// 	win := wobj.(*Window)
	// 	event.Emit("/app/finalizedWindow", event.NewValue(win))
	// 	if empty {
	// 		ievent.Emit("/app/window-all-closed", nil)
	// 	}
	// }); err != nil {
	// 	return err
	// }
	// if err := appg.AddHandler("/window/:id/changeRoute", onChangeRoute); err != nil {
	// 	return err
	// }
	// if err := appg.AddHandler("/window/:id/closed", onWindowClosed); err != nil {
	// 	return err
	// }
	// if err := appg.AddHandler("/window/:id/resize", onWindowResized); err != nil {
	// 	return err
	// }
	// if err := appg.AddHandler("/window/:id/keydown", onWindowKeyDown); err != nil {
	// 	return err
	// }
	// if err := appg.AddHandler("/window/:id/keyup", onWindowKeyUp); err != nil {
	// 	return err
	// }
	// if err := appg.AddHandler("/window/:id/focus", onWindowFocus); err != nil {
	// 	return err
	// }
	// if err := ievent.AddHandler("/window/:id/blur", onWindowBlur); err != nil {
	// 	return err
	// }

	// if err := appg.AddHandler("/window/:id/onRequestAnimationFrame", func(e *event.Event) {
	// 	w := getWindowByID(e.Params["id"])
	// 	var tick float64
	// 	if err := e.Argument.Decode(&tick); err != nil {
	// 		panic(err)
	// 	}
	// 	w.onRequestAnimationFrame(tick)
	// }); err != nil {
	// 	return err
	// }

	// if err := appg.AddHandler("/window/:id/ready", func(e *event.Event) {
	// 	w := getWindowByID(e.Params["id"])
	// 	w.onReady()
	// }); err != nil {
	// 	return err
	// }
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
	r := event.EmitWithResult(w, "new", event.NewValue(cfg))
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

func GetWindowFromBuilder(b markup.Builder) (*Window, error) {
	buildable := b.Buildable()
	if w, ok := buildable.(*Window); ok {
		return w, nil
	}
	return nil, fmt.Errorf("invalid argument: %v", buildable)
}
