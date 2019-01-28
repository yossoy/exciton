package markup

import (
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/internal/markup"
)

func ConnectToSlot(name string, sig *event.Signal) MarkupOrChild {
	return markup.SigConnecter{
		Name: name,
		Sig:  sig,
	}
}

func ConnectToSignal(name string, slot *event.Slot) MarkupOrChild {
	return markup.SlotConnecter{
		Name: name,
		Slot: slot,
	}
}
