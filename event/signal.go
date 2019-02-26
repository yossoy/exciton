package event

type Signal struct {
	slot       *Slot // TODO: remove this
	name       string
	hostTarget EventTarget
}

func (sig *Signal) Register(name string, target EventTarget) {
	sig.name = name
	sig.hostTarget = target
}

func (sig *Signal) EventPathString() string {
	if sig.hostTarget == nil {
		return ""
	}
	return EventTargetToPathString(sig.hostTarget, sig.name)
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
