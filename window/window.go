package window

import (
	"errors"
	"fmt"
	"strings"

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

func onWindowFinalize(e *event.Event) error {
	w, ok := e.Target.(*Window)
	if !ok {
		return fmt.Errorf("invalid target: %v", e.Target)
	}
	_, empty, err := object.Windows.Delete(object.ObjectKey(w.ID))
	if err != nil {
		log.PrintError("invalid window: %v", w)
		return err
	}
	// TODO: emit to AppClass?
	event.Emit(w.Parent(), "finalizeWindow", event.NewValue(w))
	if empty {
		event.Emit(w.Parent(), "window-all-closed", nil)
	}
	return nil
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
	WindowClass.AddHandler("onRequestAnimationFrame", func(e *event.Event) error {
		w, ok := e.Target.(*Window)
		if !ok {
			return fmt.Errorf("invalid target: %v", e.Target)
		}
		var tick float64
		if err := e.Argument.Decode(&tick); err != nil {
			return err
		}
		w.onRequestAnimationFrame(tick)
		return nil
	})
	WindowClass.AddHandler("ready", func(e *event.Event) error {
		w, ok := e.Target.(*Window)
		if !ok {
			return fmt.Errorf("invalid target: %v", e.Target)
		}
		w.onReady() // TODO
		return nil
	})
	markup.InitEvents(&WindowClass)

	// TODO: child event
}

// UserData is window binded data
type UserData interface{}

type Owner interface {
	event.EventTarget
	event.EventTargetWithScopedNameResolver
	//	ID() string
	PreferredLanguages() lang.PreferredLanguages
	URLBase() string
	OnActiveWindowChange(w *Window, actived bool)
}

// Window is browser window
type Window struct {
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

func (w *Window) GetTargetByScopedName(scopedName string) (event.EventTarget, string) {
	if strings.HasPrefix(scopedName, "win.") {
		return w, strings.TrimPrefix(scopedName, "win.")
	}
	return w.owner.GetTargetByScopedName(scopedName)
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

func onWindowClosed(e *event.Event) error {
	w, ok := e.Target.(*Window)
	if !ok {
		return fmt.Errorf("invalid target: %v", e.Target)
	}
	w.owner.OnActiveWindowChange(w, false)
	if w.OnClosed != nil {
		w.OnClosed(e)
	}
	return nil
}

func onWindowResized(e *event.Event) error {
	w, ok := e.Target.(*Window)
	if !ok {
		return fmt.Errorf("invalid target: %v", e.Target)
	}
	sz := geom.Size{}
	if err := e.Argument.Decode(&sz); err != nil {
		log.PrintError("/window/:id/resize: parameter decode failed: %#v", e.Argument)
		return err
	}
	log.PrintDebug("Window: resized (%#v)", sz)
	if w.OnResize != nil {
		w.OnResize(sz.Width, sz.Height)
	}
	return nil
}

func onWindowKeyDown(e *event.Event) error {
	w, ok := e.Target.(*Window)
	if !ok {
		return fmt.Errorf("invalid target: %v", e.Target)
	}
	var ke html.KeyboardEvent
	err := e.Argument.Decode(&ke)
	if err != nil {
		log.PrintError("/window/:id/keydown: parameter decode failed: %#v", e.Argument)
		return err
	}
	log.PrintDebug("Window: keydown (%#v)", ke)
	if w.OnKeyDown != nil {
		w.OnKeyDown(w, &ke)
	}
	return nil
}

func onWindowKeyUp(e *event.Event) error {
	w, ok := e.Target.(*Window)
	if !ok {
		return fmt.Errorf("invalid target: %v", e.Target)
	}
	var ke html.KeyboardEvent
	err := e.Argument.Decode(&ke)
	if err != nil {
		log.PrintError("/window/:id/keyup: parameter decode failed: %#v", e.Argument)
		return err
	}
	log.PrintDebug("Window: keyup (%#v)", ke)
	if w.OnKeyUp != nil {
		w.OnKeyUp(w, &ke)
	}
	return nil
}

func onWindowFocus(e *event.Event) error {
	w, ok := e.Target.(*Window)
	if !ok {
		return fmt.Errorf("invalid target: %v", e.Target)
	}

	w.owner.OnActiveWindowChange(w, true)
	return nil
}

func onWindowBlur(e *event.Event) error {
	w, ok := e.Target.(*Window)
	if !ok {
		return fmt.Errorf("invalid target: %v", e.Target)
	}
	w.owner.OnActiveWindowChange(w, false)
	return nil
}

type changeRouteArg struct {
	Route string `json:"route"`
}

func onChangeRoute(e *event.Event) error {
	w, ok := e.Target.(*Window)
	if !ok {
		return fmt.Errorf("invalid target: %v", e.Target)
	}
	var arg changeRouteArg
	if err := e.Argument.Decode(&arg); err != nil {
		panic(err)
	}
	log.PrintDebug("onChangeRoute: %q", arg.Route)
	w.Builder().OnRedirect(arg.Route)
	return nil
}

func InitWindows(si *driver.StartupInfo) error {
	// appg := si.AppEventRoot
	if err := initHTML(si); err != nil {
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
