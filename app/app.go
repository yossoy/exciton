package app

import (
	"fmt"

	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/internal/object"
)

type App struct {
	ID         string
	DriverData interface{}
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
