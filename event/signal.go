package event

type Signal struct {
	slot *Slot
}

func (sig *Signal) Self() *Signal {
	return sig
}

func (sig *Signal) Emit(v Value) error {
	if sig.slot != nil {
		return sig.slot.Emit(v)
	}
	return nil
}

func (sig *Signal) Disconnect() {
	if sig.slot != nil {
		sig.slot.Disconnect(sig)
	}
}
