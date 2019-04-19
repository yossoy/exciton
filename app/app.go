package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/yossoy/exciton/dialog"
	"github.com/yossoy/exciton/event"
	idialog "github.com/yossoy/exciton/internal/dialog"
	"github.com/yossoy/exciton/internal/markup"
	"github.com/yossoy/exciton/internal/object"
	"github.com/yossoy/exciton/lang"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
)

type appClass struct {
	event.EventHostCore
}

func (ac *appClass) GetTarget(id string, parent event.EventTarget) event.EventTarget {
	if ac.IsSingleton() {
		// if id != event.HostID {
		// 	panic(fmt.Sprintf("invalid id: %q", id))
		// }
		id = object.SingletonName
	}
	a := object.Apps.Get(id)
	if a == nil {
		return nil
	}
	if app, ok := a.(*App); ok {
		return app
	}
	panic(fmt.Sprintf("invalid object: %v", a))
	return nil
}

var AppClass appClass

type Owner interface {
	PreferredLanguages() lang.PreferredLanguages
}

func InitEvents(isSingleton bool, si *StartupInfo) {
	if isSingleton {
		event.InitSingletonRoot(&AppClass, "app")
	} else {
		event.InitHost(&AppClass, "app", nil)
	}
	window.InitEvents(&AppClass)
	menu.InitEvents(&AppClass)
	si.StartupInfo.AppEventHost = &AppClass
	si.StartupInfo.WinEventHost = &window.WindowClass
}

type appEventSlot struct {
	core event.Slot
}

func (a *appEventSlot) Core() *event.Slot {
	return &a.core
}

func (a *appEventSlot) Bind(h event.Handler) {
	a.core.Bind(h)
}

func (a *appEventSlot) BindWithResult(h event.HandlerWithResult) {
	a.core.BindWithResult(h)
}

func (a *appEventSlot) IsEnabled() bool {
	return a.core.IsEnabled()
}

func (a *appEventSlot) SetValidateEnabledHandler(validator func(name string) bool) {
	a.core.SetValidateEnabledHandler(validator)
}

type appEventSignal struct {
	core event.Signal
}

func (a *appEventSignal) Core() *event.Signal {
	return &a.core
}

func (a *appEventSignal) Emit(v event.Value) error {
	return a.core.Emit(v)
}

type App struct {
	owner              Owner
	id                 string
	MainWindow         *window.Window
	userDefinedSignals map[string]*appEventSignal
	userDefinedSlots   map[string]*appEventSlot
	UserData           interface{}
	activeWindow       *window.Window
}

func (app *App) RegisterEventSignal(name string) event.Signaller {
	if app.userDefinedSignals == nil {
		app.userDefinedSignals = make(map[string]*appEventSignal)
	}
	ap := &appEventSignal{}
	ap.core.Register(name, app)
	app.userDefinedSignals[name] = ap

	return ap
}

func (app *App) RegisterEventSlot(name string, handler event.Handler) event.Slotter {
	if app.userDefinedSlots == nil {
		app.userDefinedSlots = make(map[string]*appEventSlot)
	}
	s := &appEventSlot{}
	s.core.Register(name, app)
	s.core.Bind(handler)
	app.userDefinedSlots[name] = s
	return s
}

func (app *App) Parent() event.EventTarget {
	return nil
}
func (app *App) Host() event.EventHost {
	return &AppClass
}

func (app *App) GetEventSignal(name string) *event.Signal {
	if app.userDefinedSignals == nil {
		return nil
	}
	s, ok := app.userDefinedSignals[name]
	if !ok {
		return nil
	}
	return s.Core()
}

func (app *App) GetEventSlot(name string) *event.Slot {
	if app.userDefinedSlots == nil {
		return nil
	}
	s, ok := app.userDefinedSlots[name]
	if !ok {
		return nil
	}
	return s.Core()
}

func (app *App) TargetID() string {
	return app.id
}

func (app *App) ParentTarget() event.EventTarget {
	// app is root target
	return nil
}

func (app *App) GetTargetByScopedName(scopedName string) (event.EventTarget, string) {
	switch {
	case strings.HasPrefix(scopedName, "app."):
		return app, strings.TrimPrefix(scopedName, "app.")
	case strings.HasPrefix(scopedName, "win."):
		return app.activeWindow, strings.TrimPrefix(scopedName, "win.")
	default:
		return nil, scopedName
	}
}

func (app *App) PreferredLanguages() lang.PreferredLanguages {
	return app.owner.PreferredLanguages()
}

func (app *App) Owner() Owner {
	return app.owner
}

func (app *App) URLBase() string {
	if app.id == object.SingletonName {
		return ""
	}
	return "/exciton/" + app.id
}

func (app *App) OnActiveWindowChange(w *window.Window, actived bool) {
	if actived {
		app.activeWindow = w
	} else if app.activeWindow == w {
		app.activeWindow = nil
	}
}

func (app *App) ActiveWindow() *window.Window {
	return app.activeWindow
}

func (app *App) ResolveEventScopeNameToTarget(scope string) event.EventTarget {
	switch scope {
	case "app":
		return app
	case "win":
		return app.activeWindow
	default:
		// TODO: need "component" scope?
		return nil
	}
}

func NewApp(owner Owner) *App {
	id := object.Apps.NewKey()
	a := &App{
		id:    id,
		owner: owner,
	}
	object.Apps.Put(id, a)
	return a
}

func NewSingletonApp(owner Owner) *App {
	a := &App{
		id:    object.SingletonName,
		owner: owner,
	}
	object.Apps.Put(object.SingletonName, a)
	return a
}

func GetAppByID(id string) *App {
	//TODO: change to internal function?
	a := object.Apps.Get(id)
	if a == nil {
		return nil
	}
	if app, ok := a.(*App); ok {
		return app
	}
	return nil
}

func GetAppFromEventTarget(e *markup.EventTarget) (*App, error) {
	appid := e.AppID
	if appid == "" {
		appid = object.SingletonName
	}
	a := GetAppByID(appid)
	if a == nil {
		return nil, fmt.Errorf("App not found")
	}
	return a, nil
}

func GetAppFromEvent(e *event.Event) (*App, error) {
	t := e.Target
	log.Printf("Target: %v", t)
	for t != nil {
		app, ok := t.(*App)
		if ok {
			return app, nil
		}
		t = t.ParentTarget()
	}
	return nil, fmt.Errorf("App not found")
}

func (app *App) ShowMessageBoxAsync(message string, title string, messageBoxType dialog.MessageBoxType, cfg *dialog.MessageBoxConfig, handler func(int, error)) error {
	return idialog.ShowMessageBoxAsync(app, "", message, title, messageBoxType, cfg, handler)
}

func (app *App) ShowMessageBox(message string, title string, messageBoxType dialog.MessageBoxType, cfg *dialog.MessageBoxConfig) (int, error) {
	return idialog.ShowMessageBox(app, "", message, title, messageBoxType, cfg)
}

func (app *App) ShowOpenDialogAsync(cfg *dialog.FileDialogConfig, handler func(*dialog.OpenFileResult, error)) error {
	return idialog.ShowOpenDialogAsync(app, "", cfg, handler)
}

func (app *App) ShowOpenDialog(cfg *dialog.FileDialogConfig) (*dialog.OpenFileResult, error) {
	return idialog.ShowOpenDialog(app, "", cfg)
}

func (app *App) ShowSaveDialogAsync(cfg *dialog.FileDialogConfig, handler func(string, error)) error {
	return idialog.ShowSaveDialogAsync(app, "", cfg, handler)
}

func (app *App) ShowSaveDialog(cfg *dialog.FileDialogConfig) (string, error) {
	return idialog.ShowSaveDialog(app, "", cfg)
}
