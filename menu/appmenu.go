package menu

import (
	"github.com/yossoy/exciton/event"
)

func SetApplicationMenu(owner Owner, menu AppMenuTemplate) error {
	r, err := toAppMenu(menu)
	if err != nil {
		return err
	}
	mi, err := newInstance(owner, r)
	if err != nil {
		return err
	}
	return event.Emit(mi, "setApplicationMenu", event.NewValue(nil))
}
