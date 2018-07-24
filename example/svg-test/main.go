package main

import (
	"github.com/yossoy/exciton"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/log"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/svg"
	"github.com/yossoy/exciton/window"
)

type testComponent struct {
	markup.Core
	imageSVG markup.RenderResult
}

func (c *testComponent) Initialize() {
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

func (c *testComponent) Render() markup.RenderResult {
	return html.Div(
		c.imageSVG,
		html.Image(
			markup.Attribute("src", "/resources/liberty.svg"),
			markup.Attribute("width", 200),
			markup.Style("float", "right"),
		),
	)
}

var TestComponent = markup.MustRegisterComponent((*testComponent)(nil))

func onAppStart() {
	log.PrintInfo("onAppStart")
	cfg := window.WindowConfig{
		Title: "Exciton Sample",
	}
	w, err := window.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	w.Mount(TestComponent())
}

func ExcitonStartup(info *exciton.StartupInfo) error {
	info.OnAppStart = onAppStart
	info.OnAppQuit = func() {
		log.PrintInfo("app is terminated...")
	}
	if err := exciton.Init(info); err != nil {
		return err
	}
	return nil
}
