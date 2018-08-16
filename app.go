package exciton

import (
	"github.com/yossoy/exciton/event"
	ievent "github.com/yossoy/exciton/internal/event"
)

const (
	appOnQuit = "/app/finalize"
)

func Quit() {
	err := ievent.Emit("/app/quit", event.NewValue(nil))
	if err != nil {
		panic(err)
	}
}
