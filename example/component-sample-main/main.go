package main

import (
	"errors"
	"fmt"

	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/dialog"
	"github.com/yossoy/exciton/event"
	sample "github.com/yossoy/exciton/example/component-sample"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/log"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/window"
)

type rootComponent struct {
	markup.Core
	components    []string
	createdCount  int
	KillMeClicked event.Slot `exciton:"killMeClicked"`
	Clicked1      event.Slot `exciton:"clicked1"`
}

func (rc *rootComponent) Initialize() {
	rc.KillMeClicked.Bind(rc.onKillmeClicked)
	rc.Clicked1.Bind(rc.onClick1)
}

func (rc *rootComponent) newComponentKey() string {
	key := fmt.Sprintf("Create:%d", rc.createdCount)
	rc.createdCount++
	return key
}

func (rc *rootComponent) onClick1(e *event.Event) error {
	var arg sample.Click1Arg
	if err := e.Argument.Decode(&arg); err != nil {
		log.PrintDebug("arg parse error : ===> %v", err)
		return err
	}
	log.PrintDebug("arg ===> %v", arg)
	if arg.Err != nil {
		return arg.Err
	}
	win, err := window.GetWindowFromBuilder(rc.Builder())
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("Component button1 clicked!: %d", arg.Value)
	win.ShowMessageBox(msg, "sample", dialog.MessageBoxTypeInfo, nil)
	return nil
}

func (rc *rootComponent) onClickAddComponent(e *html.MouseEvent) {
	rc.components = append(rc.components, rc.newComponentKey())
	log.PrintInfo("onClickAddCompnent: %d", len(rc.components))
	rc.Builder().Rerender()
}

func (rc *rootComponent) onClickRemoveComponent(e *html.MouseEvent) {
	if len(rc.components) > 0 {
		rc.components = rc.components[0 : len(rc.components)-1]
	}
	rc.Builder().Rerender()
}

func (rc *rootComponent) onKillmeClicked(e *event.Event) error {
	var arg sample.KillMeArg
	err := e.Argument.Decode(&arg)
	if err != nil {
		return err
	}
	c := arg.Target
	if c == nil {
		return errors.New("component not found")
	}
	k, ok := c.GetProperty("objectKey")
	if ok {
		for i, ck := range rc.components {
			if ck == k {
				rc.components = append(rc.components[:i], rc.components[i+1:]...)
				rc.Builder().Rerender()
				return nil
			}
		}
	}
	return errors.New("invalid component")
}

func (rc *rootComponent) Render() markup.RenderResult {
	log.PrintDebug("Render called!: %d", len(rc.components))
	var children markup.List
	for _, k := range rc.components {
		children = append(
			children,
			markup.Keyer(
				k,
				sample.SampleComponent(
					markup.Property("objectKey", k),
					markup.ConnectToSignal(sample.SlotClick1, &rc.Clicked1),
					markup.ConnectToSignal(sample.SlotKillMe, &rc.KillMeClicked),
				),
			),
		)
	}
	return html.Div(
		html.Div(
			html.Button(
				markup.Text("+"),
				html.OnClick(rc.onClickAddComponent),
			),
			html.Button(
				markup.If(
					len(rc.components) == 0,
					markup.Attribute("disabled", true),
				),
				markup.Text("-"),
				html.OnClick(rc.onClickRemoveComponent),
			),
		),
		children,
	)
}

func onNewWindow(app *app.App, cfg *window.WindowConfig) (markup.RenderResult, error) {
	rc, err := markup.RegisterComponent((*rootComponent)(nil))
	if err != nil {
		return nil, err
	}
	cfg.Title = "Component Sample"
	return rc(), nil
}

func ExcitonStartup(info *app.StartupInfo) error {
	info.OnNewWindow = onNewWindow
	return nil
}
