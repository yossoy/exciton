package event

// TODO: change argument Value to *Event?
type SlotHandler func(Value) error

type Slot struct {
	handler     SlotHandler
	signalPaths []string
}

func (slot *Slot) Bind(h SlotHandler) {
	slot.handler = h
}

func (slot *Slot) Emit(v Value) error {
	if slot.handler != nil {
		return slot.handler(v)
	}
	return nil
}
func (slot *Slot) Connect(sig *Signal) {
	p := sig.EventPathString()
	for _, s := range slot.signalPaths {
		if s == p {
			// already connected
			return
		}
	}
	if sig.slot != nil {
		sig.slot.Disconnect(sig)
	}
	slot.signalPaths = append(slot.signalPaths, p)
	sig.slot = slot
}

func (slot *Slot) Disconnect(sig *Signal) {
	p := sig.EventPathString()
	for i, s := range slot.signalPaths {
		if s == p {
			slot.signalPaths = append(slot.signalPaths[9:i], slot.signalPaths[i+1:len(slot.signalPaths)]...)
			sig.slot = nil
			return
		}
	}
}

func (slot *Slot) DisconnectAll() {
	sp := slot.signalPaths
	slot.signalPaths = nil
	for _, s := range sp {
		t, n, err := StringToEventTarget(s)
		if err == nil {
			if st, ok := t.(EventTargetWithSignal); ok {
				sig := st.GetEventSignal(n)
				if sig != nil {
					sig.slot = nil
				}
			}
		}
	}
}
