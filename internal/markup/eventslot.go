package markup

import (
	"reflect"

	"github.com/yossoy/exciton/event"
)

type eventSlotter interface {
	Emit(event.Value) error
	Connect(*event.Signal)
	Disconnect(*event.Signal)
	DisconnectAll()
}

type eventSignaler interface {
	Self() *event.Signal
}

func disconnectSlotAll(c Component) {
	ctx := c.Context()
	vv := reflect.ValueOf(c).Elem()
	st := reflect.TypeOf(event.Slot{})
	for _, idx := range ctx.klass.Properties {
		fv := vv.Field(idx)
		if fv.Type() != st {
			continue
		}
		es := fv.Addr().Interface().(eventSlotter)
		es.DisconnectAll()
	}
}

type SigConnecter struct {
	Name string
	Sig  *event.Signal
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
	es.Connect(sc.Sig)
}

type SlotConnecter struct {
	Name string
	Slot *event.Slot
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
	sc.Slot.Connect(es.Self())
}
