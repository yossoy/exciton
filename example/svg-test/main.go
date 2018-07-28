package main

import (
	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/svg"
	"github.com/yossoy/exciton/window"
)

type rootComponent struct {
	markup.Core
	imageSVG markup.RenderResult
}

func (c *rootComponent) Initialize() {
	rf, err := driver.ResourcesFileSystem()
	if err != nil {
		panic(err)
	}
	f, err := rf.Open("/liberty.svg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	svgr, err := svg.SVGtoRenderResult(f)
	if err != nil {
		panic(err)
	}
	c.imageSVG = svgr
}

func (c *rootComponent) Render() markup.RenderResult {
	return html.Div(
		c.imageSVG,
		html.Image(
			markup.Attribute("src", "/resources/liberty.svg"),
			markup.Attribute("width", 200),
			markup.Style("float", "right"),
		),
	)
}

func onNewWindow(cfg *window.WindowConfig) (markup.RenderResult, error) {
	cfg.Title = "SVG Example"
	rc, err := markup.RegisterComponent((*rootComponent)(nil))
	if err != nil {
		return nil, err
	}
	return rc(), nil
}

func ExcitonStartup(info *app.StartupInfo) error {
	info.OnNewWindow = onNewWindow
	return nil
}
