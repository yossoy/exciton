package sample

import (
	"encoding/json"

	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/log"
	"github.com/yossoy/exciton/markup"
)

type sampleComponent struct {
	markup.Core
	ObjectKey        string       `exciton:"objectKey"`
	Mounted          bool         `exciton:"mounted"`
	OnKillMeClicked  event.Signal `exciton:"onKillMe"`
	OnButton1Clicked event.Signal `exciton:"onClick1"`
}

// Slot Names
const (
	SlotKillMe = "onKillMe"
	SlotClick1 = "onClick1"
)

func (s *sampleComponent) getClickCount() (int64, error) {
	r, err := s.CallClientFunctionSync("clientFunc1", 0)
	if err != nil {
		return 0, err
	}
	var ret json.Number
	err = json.Unmarshal(r, &ret)
	if err != nil {
		return 0, err
	}
	return ret.Int64()
}

type KillMeArg struct {
	Target markup.Component
}

func (s *sampleComponent) onKillMeClicked(e *html.MouseEvent) {
	arg := KillMeArg{
		Target: e.Target.HostComponent(),
	}
	s.OnKillMeClicked.Emit(event.NewValue(arg))
}

type Click1Arg struct {
	Err   error
	Value int64
}

func (s *sampleComponent) onButton1Clicked(e *html.MouseEvent) {
	v, err := s.getClickCount()
	log.PrintDebug("**************** onButton1Clicked: %v, %v", v, err)
	arg := Click1Arg{
		Err:   err,
		Value: v,
	}
	s.OnButton1Clicked.Emit(event.NewValue(arg))
}

func (s *sampleComponent) Render() markup.RenderResult {
	return html.Div(
		s.Classes("Sample"),
		markup.Text("Sample Component["+s.ObjectKey+"]"),
		markup.If(
			s.Mounted,
			html.Span(
				markup.Style("color", "red"),
				markup.Text(": mounted!"),
			),
		),
		html.Div(
			s.Classes("ButtonGroup"),
			html.Button(
				s.Classes("Button"),
				markup.Text("Kill ME!"),
				html.OnClick(s.onKillMeClicked),
			),
			html.Button(
				s.Classes("Button"),
				markup.Text("Call Button1Clicked handler"),
				html.OnClick(s.onButton1Clicked),
			),
			html.Button(
				s.Classes("Button"),
				markup.Text("Call onClickClient1 handler"),
				s.ClientJSEvent("click", "onClickClient1"),
			),
		),
	)
}

func onClientMount(c markup.Component, args []event.Value) {
	sc, ok := c.(*sampleComponent)
	if !ok {
		log.PrintDebug("invalid argument: %v", c)
		return
	}
	log.PrintInfo("SampleComponent on-mount called: %v", args)
	sc.Mounted = true
	c.Builder().Rerender(c)
}

func componentInit(k markup.Klass, ii markup.InitInfo) error {
	if err := ii.AddHandler("/on-mount", onClientMount); err != nil {
		return err
	}
	return nil
}

var SampleComponent = markup.MustRegisterComponent(
	(*sampleComponent)(nil),
	markup.WithComponentStyleSheet("sample.css"),
	markup.WithComponentScript("sample.js"),
	markup.WithClassInitializer(driver.InitProcTimingPostStartup, componentInit),
)
