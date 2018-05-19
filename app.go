package exciton

import "github.com/yossoy/exciton/event"

const (
	appOnQuit = "/app/finalize"
)

func Quit() {
	err := event.Emit("/app/quit", event.NewValue(nil))
	if err != nil {
		panic(err)
	}
}
