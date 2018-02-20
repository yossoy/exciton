package exciton

import "github.com/yossoy/exciton/event"

const (
	appOnQuit = "/app/finalize"
)

// App is application object
type App struct {
}

// OnQuit is called on application stop
func (a *App) OnQuit(handler func()) {
	//TODO: once
	event.AddHandler(appOnQuit, func(e *event.Event) {
		handler()
	})
}

func Quit() {
	err := event.Emit("/app/quit", event.NewValue(nil))
	if err != nil {
		panic(err)
	}
}

var app App
