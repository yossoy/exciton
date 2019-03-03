package event

type SlotHandler func(*Event) error

type EventSlotTarget interface {
	EventTarget
	EventTargetWithSlot
}

type Slot struct {
	target      EventSlotTarget
	name        string
	handler     SlotHandler
	signalPaths []string
}

func (slot *Slot) Register(name string, target EventSlotTarget) {
	slot.name = name
	slot.target = target
}

func (slot *Slot) EventPathString() string {
	if slot.target == nil {
		return ""
	}
	return EventTargetToPathString(slot.target, slot.name)
}

func (slot *Slot) Bind(h SlotHandler) {
	slot.handler = h
}

func (slot *Slot) emit(e *Event) error {
	if slot.handler != nil {
		return slot.handler(e)
	}
	return nil
}
func (slot *Slot) Connect(sig *Signal) {
	p := sig.EventPathString()
	// TODO: lock
	for _, s := range slot.signalPaths {
		if s == p {
			// already connected
			return
		}
	}
	slot.signalPaths = append(slot.signalPaths, p)
	sig.connect(slot.EventPathString())
}

func (slot *Slot) Disconnect(sigPath string) {
	// TODO: lock
	for i, s := range slot.signalPaths {
		if s == sigPath {
			slot.signalPaths = append(slot.signalPaths[:i], slot.signalPaths[i+1:]...)
			// ここではdisconnectは必要ない?
			return
		}
	}
}

func (slot *Slot) DisconnectAll() {
	slotPath := slot.EventPathString()
	sp := slot.signalPaths
	slot.signalPaths = nil
	for _, s := range sp {
		t, n, err := StringToEventTarget(s)
		if err != nil || t == nil {
			continue
		}
		st, ok := t.(EventTargetWithSignal)
		if !ok {
			continue
		}
		sig := st.GetEventSignal(n)
		if sig == nil {
			continue
		}
		sig.Disconnect(slotPath)
	}
}
