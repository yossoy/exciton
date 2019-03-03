package markup

import (
	"reflect"

	"github.com/yossoy/exciton/event"
)

func disconnectAllSignalSlot(c Component) {
	ctx := c.Context()
	vv := reflect.ValueOf(c).Elem()
	for _, idx := range ctx.klass.Signals {
		fv := vv.Field(idx)
		ss := fv.Addr().Interface().(event.Signaller)
		ss.Core().DisconnectAll()
	}
	for _, idx := range ctx.klass.Slots {
		fv := vv.Field(idx)
		ss := fv.Addr().Interface().(event.Slotter)
		ss.Core().DisconnectAll()
	}
}

type SigConnecter struct {
	Name string
	Sig  event.Signaller
}

func (sc SigConnecter) isMarkup()                              {}
func (sc SigConnecter) isMarkupOrChild()                       {}
func (sc SigConnecter) applyToNode(b Builder, n Node, on Node) { panic(false) }
func (sc SigConnecter) applyToComponent(c Component) {
	ctx := c.Context()
	idx, ok := ctx.klass.Slots[sc.Name]
	if !ok {
		panic("invalid slot name: " + sc.Name + " in " + ctx.klass.Name())
	}
	v := reflect.ValueOf(c).Elem()
	fv := v.Field(idx)
	if !fv.CanAddr() || !fv.Addr().CanInterface() {
		panic("invalid target")
	}
	es, ok := fv.Addr().Interface().(event.Slotter)
	if !ok {
		panic("invalid target")
	}
	es.Core().Connect(sc.Sig.Core())
}

type SlotConnecter struct {
	Name string
	Slot event.Slotter
}

func (sc SlotConnecter) isMarkup()                              {}
func (sc SlotConnecter) isMarkupOrChild()                       {}
func (sc SlotConnecter) applyToNode(b Builder, n Node, on Node) { panic(false) }
func (sc SlotConnecter) applyToComponent(c Component) {
	ctx := c.Context()
	idx, ok := ctx.klass.Signals[sc.Name]
	if !ok {
		panic("invalid signal name")
	}
	v := reflect.ValueOf(c).Elem()
	fv := v.Field(idx)
	if !fv.CanAddr() || !fv.Addr().CanInterface() {
		panic("invalid target")
	}
	es, ok := fv.Addr().Interface().(event.Signaller)
	if !ok {
		panic("invalid target")
	}
	sc.Slot.Core().Connect(es.Core())
}
