package main

import (
	"runtime"

	"github.com/yossoy/exciton/log"

	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
)

const (
	isDarwin = (runtime.GOOS == "darwin")
)

var appMenu = menu.AppMenuTemplate{
	{Label: menu.AppMenuLabel, Hidden: !isDarwin,
		SubMenu: menu.MenuTemplate{
			{Role: menu.RoleAbout},
			{Label: "services", Role: menu.RoleServices},
			menu.Separator,
			{Role: menu.RoleHideOthers},
			{Role: menu.RoleUnhide},
			menu.Separator,
			{Role: menu.RoleQuit},
		}},
	{Label: "File", Hidden: isDarwin,
		SubMenu: menu.MenuTemplate{
			{Role: menu.RoleQuit},
		}},
	{Label: "History",
		SubMenu: menu.MenuTemplate{
			{Role: menu.RoleGoBack},
			{Role: menu.RoleGoForward},
		}},
	{Label: "Help", Hidden: isDarwin,
		SubMenu: menu.MenuTemplate{
			{Role: menu.RoleHelp},
		}},
}

type rootComponent struct {
	markup.Core
}

func (rc *rootComponent) onClickServerRedirect(path string) {
	rc.Builder().Redirect(path)
}

func (rc *rootComponent) onChangeSelect(e *html.Event) error {
	v, err := e.Target.Node().GetProperty("value")
	if err != nil {
		return err
	}
	log.PrintDebug("redirect ==> %v", v)
	rc.Builder().Redirect(v.(string))
	return nil
}

func (rc *rootComponent) Render() markup.RenderResult {
	return html.Div(
		html.Div(
			markup.Link("/", markup.Text("[/]")),
			markup.Link("/aaa", markup.Text("[/aaa]")),
			markup.Link("/bbb", markup.Text("[/bbb]")),
			html.Button(markup.Text("[/ccc]"), markup.OnClickRedirectTo("/ccc")),
			html.Button(markup.Text("[/ddd]"), html.OnClick(func(e *html.MouseEvent) error {
				rc.onClickServerRedirect("/ddd")
				return nil
			})),
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

func onNewWindow(app *app.App, cfg *window.WindowConfig) (markup.RenderResult, error) {
	cfg.Title = "Link Sample"
	rc, err := markup.RegisterComponent((*rootComponent)(nil))
	if err != nil {
		return nil, err
	}
	return rc(), nil
}

func ExcitonStartup(info *app.StartupInfo) error {
	info.AppMenu = appMenu
	info.OnNewWindow = onNewWindow
	return nil
}
