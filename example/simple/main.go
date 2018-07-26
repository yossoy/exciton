package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/yossoy/exciton"
	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/dialog"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/log"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
)

const (
	isDarwin = (runtime.GOOS == "darwin")
)

func onOpenFile(e *html.MouseEvent) {
	log.PrintInfo("OnOpen...")
	cfg := &dialog.FileDialogConfig{
		Properties: dialog.OpenDialogForOpenFile | dialog.OpenDialogWithMultiSelections,
		Title:      "Open Files...",
	}
	files, err := dialog.ShowOpenDialog(nil, cfg)
	if err != nil {
		log.PrintError("ShowOpenFiles error: %q", err)
		return
	}

	log.PrintInfo("select files: %#v", files)

	r, err := dialog.ShowMessageBox(nil, strings.Join(files, "\n"), "Open files:", dialog.MessageBoxTypeInfo, nil)
	if err != nil {
		log.PrintError("ShowMessageBox error: %q", err)
	} else {
		log.PrintInfo("ShowMessageBox result: %d", r)
	}

}

func onSaveFile(e *html.MouseEvent) {
	log.PrintInfo("OnSave...")
	cfg := &dialog.FileDialogConfig{
		Title: "Save Files...",
	}
	file, err := dialog.ShowSaveDialog(nil, cfg)
	if err != nil {
		log.PrintError("ShowOpenFiles error: %q", err)
		return
	}

	log.PrintInfo("select file: %#v", file)

	r, err := dialog.ShowMessageBox(nil, file, "Save files:", dialog.MessageBoxTypeWarning, nil)
	if err != nil {
		log.PrintError("ShowMessageBox error: %q", err)
	} else {
		log.PrintInfo("ShowMessageBox result: %d", r)
	}

}

var appMenu = menu.AppMenuTemplate{
	{Label: menu.AppMenuLabel, Hidden: !isDarwin,
		SubMenu: menu.MenuTemplate{
			{{Role: menu.RoleAbout}},
			{{Label: "services", Role: menu.RoleServices}},
			{
				{Role: menu.RoleHideOthers},
				{Role: menu.RoleUnhide},
			},
			{{Role: menu.RoleQuit}},
		}},

	{Label: "File",
		SubMenu: menu.MenuTemplate{
			{
				{Label: "Open", Acclerator: "CommandOrControl+O", Handler: onOpenFile},
				{Label: "Save", Acclerator: "CommandOrControl+S", Handler: onSaveFile},
				{Hidden: isDarwin, Role: menu.RoleClose},
			},
			{
				{Hidden: isDarwin, Role: menu.RoleQuit},
			},
		}},
	{Label: "Edit",
		SubMenu: menu.MenuTemplate{
			{
				{Hidden: !isDarwin, Role: menu.RoleUndo},
				{Hidden: !isDarwin, Role: menu.RoleRedo},
			},
			{
				{Role: menu.RoleCut},
				{Role: menu.RoleCopy},
				{Role: menu.RolePaste},
				{Hidden: !isDarwin, Role: menu.RolePasteAndMatchStyle},
				{Role: menu.RoleDelete},
			},
			{
				{Hidden: !isDarwin, Role: menu.RoleStartSpeaking},
				{Hidden: !isDarwin, Role: menu.RoleStopSpeaking},
			},
		}},
	{Label: "Window", Role: menu.RoleWindow,
		SubMenu: menu.MenuTemplate{
			{
				{Role: menu.RoleMinimize},
				{Hidden: !isDarwin, Role: menu.RoleClose},
				{Hidden: !isDarwin, Role: menu.RoleFront},
			},
		}},
	{Label: "Help", Role: menu.RoleHelp,
		SubMenu: menu.MenuTemplate{
			{{Hidden: isDarwin, Role: menu.RoleAbout}},
		}},
}

func onClickPopupItem(e *html.MouseEvent) {
	log.PrintInfo("select Item: %#v", e.Target.ElementID)
}

var popupMenu = menu.MenuTemplate{
	{
		{Label: "Item1", Handler: onClickPopupItem},
		{Label: "Item2", Handler: onClickPopupItem},
	},
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

func (c *testChildComponent) OnContextMenu(e *html.MouseEvent) {
	w, err := window.GetWindowFromEventTarget(e.Target)
	if err != nil {
		panic(err)
	}
	err = menu.PopupMenu(popupMenu, e.ScreenPos(), w)
	if err != nil {
		panic(err)
	}
}

var TestChildComponent = markup.MustRegisterComponent((*testChildComponent)(nil))

type testComponent struct {
	markup.Core
	Text    string `exciton:"text"`
	checked bool   `exciton:"checked"`
}

func (c *testComponent) clickHandler(e *html.MouseEvent) {
	log.PrintInfo("clickHandler is called!!")
	t := e.UIEvent.Event.Target
	if t != nil {
		v, err := t.GetProperty("type")
		if err != nil {
			log.PrintError(fmt.Sprint(err))
		} else {
			vs := v.(string)
			log.PrintInfo("type = %q (%#v)", vs, v)
		}
	}
	c.Text = c.Text + "@"
	c.Context().Builder().Rerender(c)
}

func (c *testComponent) checkHandler(e *html.MouseEvent) {
	t := e.UIEvent.Event.Target
	if t != nil {
		v, err := t.GetProperty("checked")
		if err != nil {
			log.PrintError(fmt.Sprint(err))
		} else {
			log.PrintInfo("checked = %#v", v)
			c.checked = v.(bool)
		}
	}
	c.Context().Builder().Rerender(c)
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

func ExcitonStartup(info *app.StartupInfo) error {
	info.AppMenu = appMenu
	info.OnAppStart = onAppStart
	info.OnAppQuit = func() {
		log.PrintInfo("app is terminated...")
	}
	if err := exciton.Init(info); err != nil {
		return err
	}
	return nil
}
