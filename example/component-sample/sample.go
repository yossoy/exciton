package sample

import (
	"encoding/json"

	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/log"
	"github.com/yossoy/exciton/markup"
)

type SampleCompoentIF interface {
	GetClickCount() (int64, error)
}

type sampleComponent struct {
	markup.Core
	ObjectKey      string                   `exciton:"objectKey"`
	Mounted        bool                     `exciton:"mounted"`
	KillMeClicked  func(e *html.MouseEvent) `exciton:"onKillMe"`
	Button1Clicked func(e *html.MouseEvent) `exciton:"onClick1"`
}

func (s *sampleComponent) GetClickCount() (int64, error) {
	r, err := s.CallClientFunction("clientFunc1", 0)
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
				markup.If(
					s.KillMeClicked != nil,
					html.OnClick(s.KillMeClicked),
				),
			),
			html.Button(
				s.Classes("Button"),
				markup.Text("Call Button1Clicked handler"),
				markup.If(
					s.Button1Clicked != nil,
					html.OnClick(s.Button1Clicked),
				),
			),
			html.Button(
				s.Classes("Button"),
				markup.Text("Call onClickClient1 handler"),
				s.ClientJSEvent("click", "onClickClient1"),
			),
		),
	)
}

func onClientMount(c markup.Component, e *event.Event) {
	sc, ok := c.(*sampleComponent)
	if !ok {
		log.PrintDebug("invalid argument: %v", c)
		return
	}
	log.PrintDebug("argument = %v", e.Argument)
	var arg []string
	err := e.Argument.Decode(&arg)
	if err != nil {
		log.PrintDebug("argument decode failed: %v", err)
		return
	}
	log.PrintInfo("SampleComponent on-mount called: %q", arg)
	sc.Mounted = true
	c.Builder().Rerender(c)
}

func componentInit(k *markup.Klass, ii *markup.InitInfo) error {
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
