package event

type Signal struct {
	name        string
	hostTarget  EventSignalTarget
	targetPaths []string
}

type EventSignalTarget interface {
	EventTarget
	EventTargetWithSignal
}

func (sig *Signal) Register(name string, target EventSignalTarget) {
	sig.name = name
	sig.hostTarget = target
}

func (sig *Signal) connect(slotPath string) {
	// TODO: lock
	for _, s := range sig.targetPaths {
		if s == slotPath {
			// already connected
			return
		}
	}
	sig.targetPaths = append(sig.targetPaths, slotPath)
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
	tps := sig.targetPaths // need clone or lock?
	for _, tp := range tps {
		t, n, err := StringToEventTarget(tp)
		if err != nil {
			return err
		}
		if err := Emit(t, n, v); err != nil {
			return err
		}
	}
	return nil
}

func (sig *Signal) targets() []string {
	return sig.targetPaths
}

func (sig *Signal) Disconnect(slotPath string) {
	idx := -1
	for i, p := range sig.targetPaths {
		if p == slotPath {
			idx = i
			break
		}
	}
	if idx < 0 {
		return
	}
	sig.targetPaths = append(sig.targetPaths[:idx], sig.targetPaths[idx+1:]...)
	t, n, err := StringToEventTarget(slotPath)
	if err != nil {
		return
	}
	st, ok := t.(EventTargetWithSlot)
	if !ok {
		return
	}
	slot := st.GetEventSlot(n)
	if slot == nil {
		return
	}
	slot.Disconnect(sig.EventPathString())
}

func (sig *Signal) DisconnectAll() {
	sigPath := sig.EventPathString()
	tp := sig.targetPaths
	sig.targetPaths = nil
	for _, s := range tp {
		t, n, err := StringToEventTarget(s)
		if err != nil || t == nil {
			continue
		}
		st, ok := t.(EventTargetWithSlot)
		if !ok {
			continue
		}
		slot := st.GetEventSlot(n)
		if slot == nil {
			continue
		}
		slot.Disconnect(sigPath)
	}
}
