package menu

import (
	"github.com/yossoy/exciton/event"
	ievent "github.com/yossoy/exciton/internal/event"
)

func SetApplicationMenu(eventRoot string, menu AppMenuTemplate) error {
	r, err := toAppMenu(menu)
	if err != nil {
		return err
	}
	mi, err := newInstance(eventRoot, r)
	if err != nil {
		return err
	}
	return ievent.Emit(eventRoot+"/menu/"+mi.uuid+"/setApplicationMenu", event.NewValue(nil))
}
