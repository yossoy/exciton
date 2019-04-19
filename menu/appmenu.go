package menu

import "github.com/yossoy/exciton/event"

func SetApplicationMenu(owner Owner, menu AppMenuTemplate) error {
	mi, err := newAppMenu(owner, owner, menu)
	if err != nil {
		return err
	}
	return event.Emit(mi, "setApplicationMenu", event.NewValue(nil))
}
