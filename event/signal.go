package event

type signalPathItem struct {
	path  string
	names []string
}

type Signal struct {
	name        string
	hostTarget  EventSignalTarget
	targetPaths []*signalPathItem
}

type EventSignalTarget interface {
	EventTarget
	EventTargetWithSignal
}

func (sig *Signal) Register(name string, target EventSignalTarget) {
	sig.name = name
	sig.hostTarget = target
}

func (sig *Signal) connect(slotPath string, name string) {
	// TODO: lock
	for _, s := range sig.targetPaths {
		if s.path == slotPath {
			// already connected target
			for _, n := range s.names {
				if n == name {
					// alredy connect by same name
					return
				}
			}
			s.names = append(s.names, name)
			return
		}
	}
	sig.targetPaths = append(sig.targetPaths, &signalPathItem{path: slotPath, names: []string{name}})
}

func (sig *Signal) eventPathNameString() (string, string) {
	if sig.hostTarget == nil {
		return "", ""
	}
	return EventTargetToPathString(sig.hostTarget), sig.name
}

func (sig *Signal) Self() *Signal {
	return sig
}

func (sig *Signal) Emit(v Value) error {
	tps := sig.targetPaths // need clone or lock?
	for _, tp := range tps {
		t, err := StringToEventTarget(tp.path)
		if err != nil {
			return err
		}
		for _, n := range tp.names {
			b, err := IsEnableEvent(t, n)
			if err != nil {
				return err
			}
			if b {
				if err := Emit(t, n, v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (sig *Signal) IsEnabled() bool {
	tps := sig.targetPaths // need clone or lock?
	enabled := false
	for _, tp := range tps {
		t, err := StringToEventTarget(tp.path)
		if err != nil {
			continue
		}
		for _, n := range tp.names {
			b, err := IsEnableEvent(t, n)
			if err == nil {
				enabled = enabled || b
			}
		}
	}
	return enabled
}

func (sig *Signal) targets() []*signalPathItem {
	return sig.targetPaths
}

func (sig *Signal) Disconnect(slotPath string, name string) {
	pathIdx := -1
	for i, p := range sig.targetPaths {
		if p.path == slotPath {
			pathIdx = i
			break
		}
	}
	if pathIdx < 0 {
		return
	}
	sp := sig.targetPaths[pathIdx]
	nameIdx := -1
	for i, n := range sp.names {
		if n == name {
			nameIdx = i
			break
		}
	}
	if nameIdx < 0 {
		return
	}
	if len(sp.names) == 1 {
		sig.targetPaths = append(sig.targetPaths[:pathIdx], sig.targetPaths[pathIdx+1:]...)
	} else {
		sp.names = append(sp.names[:nameIdx], sp.names[nameIdx+1:]...)
	}

	t, err := StringToEventTarget(slotPath)
	if err != nil {
		return
	}
	st, ok := t.(EventTargetWithSlot)
	if !ok {
		return
	}
	slot := st.GetEventSlot(name)
	if slot == nil {
		return
	}
	slot.Disconnect(sig.eventPathNameString())
}

func (sig *Signal) DisconnectAll() {
	sigPath, name := sig.eventPathNameString()
	tp := sig.targetPaths
	sig.targetPaths = nil
	for _, s := range tp {
		t, err := StringToEventTarget(s.path)
		if err != nil || t == nil {
			continue
		}
		st, ok := t.(EventTargetWithSlot)
		if !ok {
			continue
		}
		for _, n := range s.names {
			slot := st.GetEventSlot(n)
			if slot == nil {
				continue
			}
			slot.Disconnect(sigPath, name)
		}
	}
}
