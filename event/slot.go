package event

type EventSlotTarget interface {
	EventTarget
	EventTargetWithSlot
}

type slotPathItem struct {
	path  string
	names []string
}

type Slot struct {
	target            EventSlotTarget
	name              string
	handler           Handler
	handlerWithResult HandlerWithResult
	signalPaths       []*slotPathItem
	validator         func(name string) bool
}

func (slot *Slot) Register(name string, target EventSlotTarget) {
	slot.name = name
	slot.target = target
}

func (slot *Slot) eventPathNameString() (string, string) {
	if slot.target == nil {
		return "", ""
	}
	return EventTargetToPathString(slot.target), slot.name
}

func (slot *Slot) Bind(h Handler) {
	slot.handler = h
	slot.handlerWithResult = nil
}

func (slot *Slot) BindWithResult(h HandlerWithResult) {
	slot.handler = nil
	slot.handlerWithResult = h
}

func (slot *Slot) IsEnabled() bool {
	if slot.validator != nil {
		return slot.validator(slot.name)
	}
	return true
}

func (slot *Slot) SetValidateEnabledHandler(validator func(name string) bool) {
	slot.validator = validator
}

func (slot *Slot) emit(e *Event) error {
	if slot.handler != nil && slot.IsEnabled() {
		return slot.handler(e)
	}
	return nil
}
func (slot *Slot) Connect(sig *Signal) {
	sigPath, sigName := sig.eventPathNameString()
	// TODO: lock
	for _, s := range slot.signalPaths {
		if s.path == sigPath {
			for _, n := range s.names {
				if n == sigName {
					// already connect same name
					return
				}
			}
			s.names = append(s.names, sigName)
			sig.connect(slot.eventPathNameString())
			return
		}
	}
	slot.signalPaths = append(slot.signalPaths, &slotPathItem{path: sigPath, names: []string{sigName}})
	sig.connect(slot.eventPathNameString())
}

func (slot *Slot) Disconnect(sigPath string, name string) {
	// TODO: lock
	for i, s := range slot.signalPaths {
		if s.path == sigPath {
			for j, n := range s.names {
				if n == name {
					if len(s.names) == 1 {
						slot.signalPaths = append(slot.signalPaths[:i], slot.signalPaths[i+1:]...)
					} else {
						s.names = append(s.names[:j], s.names[j+1:]...)
					}
					return
				}
			}
		}
	}
}

func (slot *Slot) DisconnectAll() {
	slotPath, slotName := slot.eventPathNameString()
	sp := slot.signalPaths
	slot.signalPaths = nil
	for _, s := range sp {
		t, err := StringToEventTarget(s.path)
		if err != nil || t == nil {
			continue
		}
		st, ok := t.(EventTargetWithSignal)
		if !ok {
			continue
		}
		for _, n := range s.names {
			sig := st.GetEventSignal(n)
			if sig == nil {
				continue
			}
			sig.Disconnect(slotPath, slotName)
		}
	}
}
