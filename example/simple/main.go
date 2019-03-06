package main

import (
	"runtime"
	"strings"

	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/dialog"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/log"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
)

func onOpenFile(e *html.MouseEvent) error {
	app, err := app.GetAppFromEventTarget(e.Target)
	if err != nil {
		return err
	}
	log.PrintInfo("OnOpen...")
	cfg := &dialog.FileDialogConfig{
		Properties: dialog.OpenDialogForOpenFile | dialog.OpenDialogWithMultiSelections,
		Title:      "Open Files...",
	}
	files, err := app.ShowOpenDialog(cfg)
	if err != nil {
		log.PrintError("ShowOpenFiles error: %q", err)
		return err
	}
	defer files.Cleanup()

	log.PrintInfo("select files: %#v", files)
	fileNames := []string{}
	for _, f := range files.Items {
		fileNames = append(fileNames, f.Name())
	}

	r, err := app.ShowMessageBox(strings.Join(fileNames, "\n"), "Open files:", dialog.MessageBoxTypeInfo, nil)
	if err != nil {
		log.PrintError("ShowMessageBox error: %q", err)
	} else {
		log.PrintInfo("ShowMessageBox result: %d", r)
	}
	return err
}

func onSaveFile(e *html.MouseEvent) error {
	app, err := app.GetAppFromEventTarget(e.Target)
	if err != nil {
		return err
	}
	log.PrintInfo("OnSave...")
	cfg := &dialog.FileDialogConfig{
		Title: "Save Files...",
	}
	file, err := app.ShowSaveDialog(cfg)
	if err != nil {
		log.PrintError("ShowOpenFiles error: %q", err)
		return err
	}

	log.PrintInfo("select file: %#v", file)

	r, err := app.ShowMessageBox(file, "Save files:", dialog.MessageBoxTypeWarning, nil)
	if err != nil {
		log.PrintError("ShowMessageBox error: %q", err)
	} else {
		log.PrintInfo("ShowMessageBox result: %d", r)
	}
	return err
}

func appMenu(isDarwin bool) menu.AppMenuTemplate {
	return menu.AppMenuTemplate{
		{Label: menu.AppMenuLabel, Hidden: !isDarwin,
			SubMenu: menu.MenuTemplate{
				{Role: menu.RoleAbout},
				menu.Separator,
				{Label: "services", Role: menu.RoleServices},
				menu.Separator,
				{Role: menu.RoleHideOthers},
				{Role: menu.RoleUnhide},
				menu.Separator,
				{Role: menu.RoleQuit},
			}},

		{Label: "File",
			SubMenu: menu.MenuTemplate{
				{Label: "Open", Acclerator: "CommandOrControl+O", Handler: onOpenFile},
				{Label: "Save", Acclerator: "CommandOrControl+S", Handler: onSaveFile},
				{Hidden: isDarwin, Role: menu.RoleClose},
				menu.Separator,
				{Hidden: isDarwin, Role: menu.RoleQuit},
			}},
		{Label: "Edit",
			SubMenu: menu.MenuTemplate{
				{Hidden: !isDarwin, Role: menu.RoleUndo},
				{Hidden: !isDarwin, Role: menu.RoleRedo},
				menu.Separator,
				{Role: menu.RoleCut},
				{Role: menu.RoleCopy},
				{Role: menu.RolePaste},
				{Hidden: !isDarwin, Role: menu.RolePasteAndMatchStyle},
				{Role: menu.RoleDelete},
				menu.Separator,
				{Hidden: !isDarwin, Role: menu.RoleStartSpeaking},
				{Hidden: !isDarwin, Role: menu.RoleStopSpeaking},
			}},
		{Label: "Window", Role: menu.RoleWindow,
			SubMenu: menu.MenuTemplate{
				{Role: menu.RoleMinimize},
				{Hidden: !isDarwin, Role: menu.RoleClose},
				{Hidden: !isDarwin, Role: menu.RoleFront},
			}},
		{Label: "Help", Role: menu.RoleHelp,
			SubMenu: menu.MenuTemplate{
				{Hidden: isDarwin, Role: menu.RoleAbout},
			}},
	}
}

func onClickPopupItem(e *html.MouseEvent) error {
	log.PrintInfo("select Item: %#v", e.Target.ElementID)
	return nil
}

var popupMenu = menu.MenuTemplate{
	{Label: "Item1", Handler: onClickPopupItem},
	{Label: "Item2", Handler: onClickPopupItem},
}

type testChildComponent struct {
	markup.Core
	Text string `exciton:"text"`
}

func (c *testChildComponent) Render() markup.RenderResult {
	return html.Span(
		markup.Style("color", "red"),
		markup.Style("background-color", "green"),
		html.ContextMenu(c.OnContextMenu).PreventDefault().StopPropagation(),
		markup.Text(c.Text),
	)
}

func (c *testChildComponent) OnContextMenu(e *html.MouseEvent) error {
	w, err := window.GetWindowFromEventTarget(e.Target)
	if err != nil {
		return err
	}
	err = menu.PopupMenu(popupMenu, e.ScreenPos(), w)
	if err != nil {
		return err
	}
	return nil
}

var TestChildComponent = markup.MustRegisterComponent((*testChildComponent)(nil))

type testComponent struct {
	markup.Core
	Text    string `exciton:"text"`
	checked bool   `exciton:"checked"`
}

func (c *testComponent) clickSlotHandler(e *event.Event) error {
	log.PrintDebug("component clickSlotHandler is called: %v", *e)
	return nil
}

func (c *testComponent) clickHandler(e *html.MouseEvent) error {
	log.PrintInfo("clickHandler is called!!")
	n := e.UIEvent.Event.Target.Node()
	if n != nil {
		v, err := n.GetProperty("type")
		if err != nil {
			log.PrintError("%v", err)
			return err
		} else {
			vs := v.(string)
			log.PrintInfo("type = %q (%#v)", vs, v)
		}
	}
	c.Text = c.Text + "@"
	c.Context().Builder().Rerender(c)
	return nil
}

func (c *testComponent) checkHandler(e *html.MouseEvent) error {
	n := e.UIEvent.Event.Target.Node()
	if n != nil {
		v, err := n.GetProperty("checked")
		if err != nil {
			log.PrintError("%v", err)
			return err
		} else {
			log.PrintInfo("checked = %#v", v)
			c.checked = v.(bool)
		}
	}
	c.Context().Builder().Rerender(c)
	return nil
}

func (c *testComponent) Render() markup.RenderResult {
	return html.Div(
		html.Image(
			markup.Attribute("src", "/resources/liberty.svg"),
			markup.Attribute("width", 200),
			markup.Style("float", "right"),
		),
		html.Heading1(
			markup.Style("color", "red"),
			markup.Text("Exciton Sample"),
		),
		html.Div(
			markup.Text("dynamic added text==>"),
			markup.Text(c.Text),
			html.Button(
				html.OnClick(c.clickHandler),
				markup.Text("Click Me!"),
			),
		),
		html.Div(
			html.Label(
				markup.Attribute("for", "block_change"),
				html.Input(
					markup.Attribute("id", "block_change"),
					markup.Attribute("type", "checkbox"),
					markup.If(
						c.checked,
						markup.Attribute("checked", "checked"),
					),
					html.OnClick(c.checkHandler),
				),
				markup.Text("Show hidden contents"),
			),
		),
		markup.If(
			c.checked,
			TestChildComponent(
				markup.Property("text", "Child Component Test"),
			),
		),
	)
}

var TestComponent = markup.MustRegisterComponent((*testComponent)(nil))

func onNewWindow(app *app.App, cfg *window.WindowConfig) (markup.RenderResult, error) {
	cfg.Title = "Exciton Sample"
	return TestComponent(), nil
}

func onAppStart(app *app.App, info *app.StartupInfo) error {
	log.PrintInfo("onAppStart")
	return nil
}

func ExcitonStartup(info *app.StartupInfo) error {
	isDarwinApp := (runtime.GOOS == "darwin") && (driver.Type() != "web")
	info.AppMenu = appMenu(isDarwinApp)
	info.OnAppStart = onAppStart
	info.OnNewWindow = onNewWindow
	info.OnAppQuit = func() {
		log.PrintInfo("app is terminated...")
	}
	return nil
}
