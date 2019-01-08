package app

import (
	"fmt"

	"github.com/yossoy/exciton/dialog"
	"github.com/yossoy/exciton/event"
	idialog "github.com/yossoy/exciton/internal/dialog"
	"github.com/yossoy/exciton/internal/object"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/window"
)

type App struct {
	ID         string
	DriverData interface{}
	MainWindow *window.Window
	UserData   interface{}
}

func (app *App) eventRoot() string {
	if app.ID == object.SingletonName {
		return ""
	}
	return "/exciton/" + app.ID
}

func NewApp(driverData interface{}) *App {
	id := object.Apps.NewKey()
	a := &App{
		ID:         id,
		DriverData: driverData,
	}
	object.Apps.Put(id, a)
	return a
}

func NewSingletonApp(driverData interface{}) *App {
	a := &App{
		ID:         object.SingletonName,
		DriverData: driverData,
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
	return idialog.ShowMessageBoxAsync(app.eventRoot(), "", message, title, messageBoxType, cfg, handler)
}

func (app *App) ShowMessageBox(message string, title string, messageBoxType dialog.MessageBoxType, cfg *dialog.MessageBoxConfig) (int, error) {
	return idialog.ShowMessageBox(app.eventRoot(), "", message, title, messageBoxType, cfg)
}

func (app *App) ShowOpenDialogAsync(cfg *dialog.FileDialogConfig, handler func(*dialog.OpenFileResult, error)) error {
	return idialog.ShowOpenDialogAsync(app.eventRoot(), "", cfg, handler)
}

func (app *App) ShowOpenDialog(cfg *dialog.FileDialogConfig) (*dialog.OpenFileResult, error) {
	return idialog.ShowOpenDialog(app.eventRoot(), "", cfg)
}

func (app *App) ShowSaveDialogAsync(cfg *dialog.FileDialogConfig, handler func(string, error)) error {
	return idialog.ShowSaveDialogAsync(app.eventRoot(), "", cfg, handler)
}

func (app *App) ShowSaveDialog(cfg *dialog.FileDialogConfig) (string, error) {
	return idialog.ShowSaveDialog(app.eventRoot(), "", cfg)
}
