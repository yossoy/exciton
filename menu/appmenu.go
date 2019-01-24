package menu

import (
	"github.com/yossoy/exciton/event"
	ievent "github.com/yossoy/exciton/internal/event"
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
	return ievent.Emit(mi.EventPath("setApplicationMenu"), event.NewValue(nil))
}
