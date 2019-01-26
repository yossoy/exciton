package markup

import (
	"github.com/yossoy/exciton/internal/markup"
)

type EventSlot struct {
	markup.EventSlot
}

type EventSignal struct {
	markup.EventSignal
}

func ConnectToSlot(name string, sig *EventSignal) MarkupOrChild {
	return markup.SigConnecter{
		Name: name,
		Sig:  &sig.EventSignal,
	}
}

func ConnectToSignal(name string, slot *EventSlot) MarkupOrChild {
	return markup.SlotConnecter{
		Name: name,
		Slot: &slot.EventSlot,
	}
}
