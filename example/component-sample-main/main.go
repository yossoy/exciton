package main

import (
	"fmt"

	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/dialog"
	"github.com/yossoy/exciton/example/component-sample"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/log"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/window"
)

type rootComponent struct {
	markup.Core
	components   []string
	createdCount int
}

func (rc *rootComponent) newComponentKey() string {
	key := fmt.Sprintf("Create:%d", rc.createdCount)
	rc.createdCount++
	return key
}

func (rc *rootComponent) onClickButton1(e *html.MouseEvent) {
	win, err := window.GetWindowFromEventTarget(e.Target)
	if err != nil {
		panic(err)
	}
	c := e.Target.HostComponent()
	if c == nil {
		panic("component not found")
	}
	sc, ok := c.(sample.SampleCompoentIF)
	if !ok {
		panic(fmt.Sprintf("invalid component: %v", c))
	}
	cc, err := sc.GetClickCount()
	if err != nil {
		panic(err)
	}
	msg := fmt.Sprintf("Component button1 clicked!: %d", cc)
	win.ShowMessageBox(msg, "sample", dialog.MessageBoxTypeInfo, nil)
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

func (rc *rootComponent) onKillmeClicked(e *html.MouseEvent) {
	c := e.Target.HostComponent()
	if c == nil {
		panic("component not found!")
	}
	k, ok := c.GetProperty("objectKey")
	if !ok {
		panic("invalid component")
	}
	for i, ck := range rc.components {
		if ck == k {
			rc.components = append(rc.components[:i], rc.components[i+1:]...)
			rc.Builder().Rerender()
			return
		}
	}
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
					markup.Property("onClick1", rc.onClickButton1),
					markup.Property("onKillMe", rc.onKillmeClicked),
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
