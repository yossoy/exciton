package markup

import (
	"reflect"

	"github.com/yossoy/exciton/event"
)

//TODO: change slot/signal to event base implement (not function pointer).

type eventSlotter interface {
	emit(event.Value) error
	connect(*EventSignal)
	disconnect()
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

type EventSlotHandler func(event.Value) error

type EventSlot struct {
	handler EventSlotHandler
	signal  *EventSignal
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
	if slot.signal != sig {
		slot.disconnect()
	}
	slot.signal = sig
	sig.slot = slot
}

func (slot *EventSlot) disconnect() {
	if slot.signal != nil {
		slot.signal.slot = nil
		slot.signal = nil
	}
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
		es.disconnect()
	}
}

type sigConnecter struct {
	name string
	sig  *EventSignal
}

func (sc sigConnecter) isMarkup()                              {}
func (sc sigConnecter) isMarkupOrChild()                       {}
func (sc sigConnecter) applyToNode(b Builder, n Node, on Node) { panic(false) }
func (sc sigConnecter) applyToComponent(c Component) {
	ctx := c.Context()
	idx, ok := ctx.klass.Properties[sc.name]
	if !ok {
		panic("invalid signal name")
	}
	v := reflect.ValueOf(c).Elem()
	fv := v.Field(idx)
	if fv.Type() != reflect.TypeOf(EventSlot{}) {
		panic("invalid target")
	}
	es := fv.Addr().Interface().(eventSlotter)
	es.connect(sc.sig)
}

func ConnectToSlot(name string, sig *EventSignal) MarkupOrChild {
	return sigConnecter{
		name: name,
		sig:  sig,
	}
}

type slotConnecter struct {
	name string
	slot *EventSlot
}

func (sc slotConnecter) isMarkup()                              {}
func (sc slotConnecter) isMarkupOrChild()                       {}
func (sc slotConnecter) applyToNode(b Builder, n Node, on Node) { panic(false) }
func (sc slotConnecter) applyToComponent(c Component) {
	ctx := c.Context()
	idx, ok := ctx.klass.Properties[sc.name]
	if !ok {
		panic("invalid signal name")
	}
	v := reflect.ValueOf(c).Elem()
	fv := v.Field(idx)
	if fv.Type() != reflect.TypeOf(EventSignal{}) {
		panic("invalid target")
	}
	es := fv.Addr().Interface().(eventSignaler)
	sc.slot.connect(es.ptr())
}

func ConnectToSignal(name string, slot *EventSlot) MarkupOrChild {
	return slotConnecter{
		name: name,
		slot: slot,
	}
}
