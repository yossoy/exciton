package markup

import (
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/internal/markup"
)

type EventSlot struct {
	core event.Slot
}

func (es *EventSlot) Core() *event.Slot {
	return &es.core
}

func (es *EventSlot) Bind(handler event.Handler) {
	es.core.Bind(handler)
}

type EventSignal struct {
	core event.Signal
}

func (es *EventSignal) Emit(v event.Value) error {
	return es.core.Emit(v)
}

func (es *EventSignal) Core() *event.Signal {
	return &es.core
}

func ConnectToSlot(name string, sig event.Signaller) MarkupOrChild {
	return markup.SigConnecter{
		Name: name,
		Sig:  sig,
	}
}

func ConnectToSignal(name string, slot event.Slotter) MarkupOrChild {
	return markup.SlotConnecter{
		Name: name,
		Slot: slot,
	}
}
