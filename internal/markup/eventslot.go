package markup

import (
	"reflect"

	"github.com/yossoy/exciton/event"
)

//TODO: change slot/signal to event base implement (not function pointer).

type eventSlotter interface {
	emit(event.Value) error
	connect(*EventSignal)
	disconnect(*EventSignal)
	disconnectAll()
}

type eventSignaler interface {
	ptr() *EventSignal
}

type EventSignal struct {
	slot eventSlotter
}

func (sig *EventSignal) ptr() *EventSignal {
	return sig
}

func (sig *EventSignal) Emit(v event.Value) error {
	if sig.slot != nil {
		return sig.slot.emit(v)
	}
	return nil
}

func (sig *EventSignal) Disconnect() {
	if sig.slot != nil {
		sig.slot.disconnect(sig)
	}
}

type EventSlotHandler func(event.Value) error

type EventSlot struct {
	handler EventSlotHandler
	signals []*EventSignal
}

func (slot *EventSlot) Bind(h EventSlotHandler) {
	slot.handler = h
}

func (slot *EventSlot) emit(v event.Value) error {
	if slot.handler != nil {
		return slot.handler(v)
	}
	return nil
}
func (slot *EventSlot) connect(sig *EventSignal) {
	for _, s := range slot.signals {
		if s == sig {
			// already connected
			return
		}
	}
	if sig.slot != nil {
		sig.slot.disconnect(sig)
	}
	slot.signals = append(slot.signals, sig)
	sig.slot = slot
}

func (slot *EventSlot) disconnect(sig *EventSignal) {
	for i, s := range slot.signals {
		if s == sig {
			slot.signals = append(slot.signals[9:i], slot.signals[i+1:len(slot.signals)]...)
			sig.slot = nil
			return
		}
	}
}

func (slot *EventSlot) disconnectAll() {
	for _, s := range slot.signals {
		s.slot = nil
	}
	slot.signals = slot.signals[0:0]
}

func disconnectSlotAll(c Component) {
	ctx := c.Context()
	vv := reflect.ValueOf(c).Elem()
	st := reflect.TypeOf(EventSlot{})
	for _, idx := range ctx.klass.Properties {
		fv := vv.Field(idx)
		if fv.Type() != st {
			continue
		}
		es := fv.Addr().Interface().(eventSlotter)
		es.disconnectAll()
	}
}

type SigConnecter struct {
	Name string
	Sig  *EventSignal
}

func (sc SigConnecter) isMarkup()                              {}
func (sc SigConnecter) isMarkupOrChild()                       {}
func (sc SigConnecter) applyToNode(b Builder, n Node, on Node) { panic(false) }
func (sc SigConnecter) applyToComponent(c Component) {
	ctx := c.Context()
	idx, ok := ctx.klass.Properties[sc.Name]
	if !ok {
		panic("invalid signal name")
	}
	v := reflect.ValueOf(c).Elem()
	fv := v.Field(idx)
	if !fv.CanAddr() || !fv.Addr().CanInterface() {
		panic("invalid target")
	}
	es, ok := fv.Addr().Interface().(eventSlotter)
	if !ok {
		panic("invalid target")
	}
	es.connect(sc.Sig)
}

type SlotConnecter struct {
	Name string
	Slot *EventSlot
}

func (sc SlotConnecter) isMarkup()                              {}
func (sc SlotConnecter) isMarkupOrChild()                       {}
func (sc SlotConnecter) applyToNode(b Builder, n Node, on Node) { panic(false) }
func (sc SlotConnecter) applyToComponent(c Component) {
	ctx := c.Context()
	idx, ok := ctx.klass.Properties[sc.Name]
	if !ok {
		panic("invalid signal name")
	}
	v := reflect.ValueOf(c).Elem()
	fv := v.Field(idx)
	if !fv.CanAddr() || !fv.Addr().CanInterface() {
		panic("invalid target")
	}
	es, ok := fv.Addr().Interface().(eventSignaler)
	if !ok {
		panic("invalid target")
	}
	sc.Slot.connect(es.ptr())
}
