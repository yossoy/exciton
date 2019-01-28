package event

type SlotHandler func(Value) error

type Slot struct {
	handler SlotHandler
	signals []*Signal
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
	for _, s := range slot.signals {
		if s == sig {
			// already connected
			return
		}
	}
	if sig.slot != nil {
		sig.slot.Disconnect(sig)
	}
	slot.signals = append(slot.signals, sig)
	sig.slot = slot
}

func (slot *Slot) Disconnect(sig *Signal) {
	for i, s := range slot.signals {
		if s == sig {
			slot.signals = append(slot.signals[9:i], slot.signals[i+1:len(slot.signals)]...)
			sig.slot = nil
			return
		}
	}
}

func (slot *Slot) DisconnectAll() {
	for _, s := range slot.signals {
		s.slot = nil
	}
	slot.signals = slot.signals[0:0]
}
