package main

import (
	"runtime"

	"github.com/yossoy/exciton"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
)

const (
	isDarwin = (runtime.GOOS == "darwin")
)

type menuBar struct {
	markup.Core
}

func (m *menuBar) Render() markup.RenderResult {
	return menu.ApplicationMenu(
		markup.If(
			isDarwin,
			menu.SubMenu("*Appname*",
				menu.RoledItem(menu.RoleAbout),
				menu.Separator(),
				menu.RoledMenu(menu.RoleServices, "services"),
				menu.Separator(),
				menu.RoledItem(menu.RoleHideOthers),
				menu.RoledItem(menu.RoleUnhide),
				menu.Separator(),
				menu.RoledItem(menu.RoleQuit),
			),
		),
		markup.If(
			!isDarwin,
			menu.SubMenu("File",
				menu.RoledItem(menu.RoleQuit),
			),
		),
		menu.SubMenu("History",
			menu.RoledItem(menu.RoleGoBack),
			menu.RoledItem(menu.RoleGoForward),
		),
		menu.RoledMenu(menu.RoleHelp, "Help",
			markup.If(
				!isDarwin,
				menu.RoledItem(menu.RoleAbout),
			),
		),
	)
}

var MenuBar = markup.MustRegisterComponent((*menuBar)(nil))

type rootComponent struct {
	markup.Core
}

func (rc *rootComponent) onClickServerRedirect(path string) {
	rc.Builder().Redirect(path)
}

func (rc *rootComponent) onChangeSelect(e *html.Event) {
	v, err := e.Target.GetProperty("value")
	if err != nil {
		panic(err)
	}
	rc.Builder().Redirect(v.(string))
}

func (rc *rootComponent) Render() markup.RenderResult {
	return html.Div(
		html.Div(
			markup.Link("/", markup.Text("[/]")),
			markup.Link("/aaa", markup.Text("[/aaa]")),
			markup.Link("/bbb", markup.Text("[/bbb]")),
			html.Button(markup.Text("[/ccc]"), markup.OnClickRedirectTo("/ccc")),
			html.Button(markup.Text("[/ddd]"), html.OnClick(func(e *html.MouseEvent) { rc.onClickServerRedirect("/ddd") })),
			html.Select(
				html.OnChange(rc.onChangeSelect),
				html.Option(markup.Attribute("value", "/eee/1"), markup.Text("[/eee/1]")),
				html.Option(markup.Attribute("value", "/eee/2"), markup.Text("[/eee/2]")),
				html.Option(markup.Attribute("value", "/eee/3"), markup.Text("[/eee/3]")),
				html.Option(markup.Attribute("value", "/eee/4"), markup.Text("[/eee/4]")),
			),
		),
		markup.BrowserRouter(
			markup.ExactRoute("/", markup.Text("Root")),
			markup.Route("/aaa", markup.Text("aaa")),
			markup.Route("/bbb", markup.Text("bbb")),
			markup.Route("/ccc", markup.Text("ccc")),
			markup.Route("/ddd", markup.Text("ddd")),
			markup.Route("/eee/:val", markup.Text("eee")),
			markup.FallbackRoute(markup.Text("invalid route!")),
		),
	)
}

func onAppStart() {
	rc := markup.MustRegisterComponent((*rootComponent)(nil))
	menu.SetApplicationMenu(menu.MustNew(MenuBar()))
	cfg := window.WindowConfig{
		Title: "Link Sample",
	}
	w, err := window.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	w.Mount(rc())
}

func ExcitonStartup(info *exciton.StartupInfo) error {
	info.OnAppStart = onAppStart
	info.OnAppQuit = func() {}
	if err := exciton.Init(info); err != nil {
		return err
	}
	return nil
}
