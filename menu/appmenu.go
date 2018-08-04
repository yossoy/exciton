package menu

import (
	"github.com/yossoy/exciton/event"
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
	return event.Emit(eventRoot+"/menu/"+mi.uuid+"/setApplicationMenu", event.NewValue(nil))
}
