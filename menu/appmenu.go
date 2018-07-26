package menu

import (
	"github.com/yossoy/exciton/event"
)

func SetApplicationMenu(menu AppMenuTemplate) error {
	r, err := toAppMenu(menu)
	if err != nil {
		return err
	}
	mi, err := newInstance(r)
	if err != nil {
		return err
	}
	return event.Emit("/menu/"+mi.uuid+"/setApplicationMenu", event.NewValue(nil))
}
