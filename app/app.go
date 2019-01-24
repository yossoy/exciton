package app

import (
	"fmt"
	"strings"

	"github.com/yossoy/exciton/dialog"
	"github.com/yossoy/exciton/event"
	idialog "github.com/yossoy/exciton/internal/dialog"
	"github.com/yossoy/exciton/internal/object"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/window"
)

type Owner interface {
	PreferredLanguages() []string
}

type App struct {
	owner      Owner
	id         string
	MainWindow *window.Window
	UserData   interface{}
}

func (app *App) ID() string {
	return app.id
}

func (app *App) PreferredLanguages() []string {
	return app.owner.PreferredLanguages()
}

func (app *App) Owner() Owner {
	return app.owner
}

func (app *App) EventPath(fragments ...string) string {
	var sb strings.Builder
	sb.Grow(100)
	if app.id != object.SingletonName {
		sb.WriteString("/exciton/")
		sb.WriteString(app.id)
	}
	for _, f := range fragments {
		sb.WriteString("/")
		sb.WriteString(f)
	}
	return sb.String()
}

func (app *App) EventPath2(fragments1 []string, fragments2 []string) string {
	var sb strings.Builder
	sb.Grow(100)
	if app.id != object.SingletonName {
		sb.WriteString("/exciton/")
		sb.WriteString(app.id)
	}
	for _, f := range fragments1 {
		sb.WriteString("/")
		sb.WriteString(f)
	}
	for _, f := range fragments2 {
		sb.WriteString("/")
		sb.WriteString(f)
	}
	return sb.String()
}

func (app *App) AppEventPath(fragments ...string) string {
	return app.EventPath(fragments...)
}

func (app *App) URLBase() string {
	if app.id == object.SingletonName {
		return ""
	}
	return "/exciton/" + app.id
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
	appid, ok := e.Params["appid"]
	if !ok {
		appid = object.SingletonName
	}
	a := GetAppByID(appid)
	if a == nil {
		return nil, fmt.Errorf("App not found")
	}
	return a, nil
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
