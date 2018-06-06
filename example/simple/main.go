package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/yossoy/exciton/menu"

	"github.com/yossoy/exciton/html"

	"github.com/yossoy/exciton"
	"github.com/yossoy/exciton/dialog"
	"github.com/yossoy/exciton/log"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/window"
)

const (
	isDarwin = (runtime.GOOS == "darwin")
)

type menuBar struct {
	markup.Core
}

func (m *menuBar) OnOpen(e *html.MouseEvent) {
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

func (m *menuBar) OnSave(e *html.MouseEvent) {
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
		menu.SubMenu("File",
			menu.Item("Open;CommandOrControl+O", m.OnOpen),
			menu.Item("Save;CommandOrControl+S", m.OnSave),
		),
		menu.SubMenu("Edit",
			markup.If(
				isDarwin,
				menu.RoledItem(menu.RoleUndo),
				menu.RoledItem(menu.RoleRedo),
				menu.Separator(),
			),
			menu.RoledItem(menu.RoleCut),
			menu.RoledItem(menu.RoleCopy),
			menu.RoledItem(menu.RolePaste),
			markup.If(
				isDarwin,
				menu.RoledItem(menu.RolePasteAndMatchStyle),
			),
			menu.RoledItem(menu.RoleDelete),
			markup.If(
				isDarwin,
				menu.Separator(),
				menu.RoledItem(menu.RoleStartSpeaking),
				menu.RoledItem(menu.RoleStopSpeaking),
			),
		),
		menu.RoledMenu(menu.RoleWindow, "Window",
			menu.RoledItem(menu.RoleMinimize),
			menu.RoledItem(menu.RoleClose),
			markup.If(
				isDarwin,
				menu.Separator(),
				menu.RoledItem(menu.RoleFront),
			),
		),
		menu.RoledMenu(menu.RoleHelp, "Help",
			markup.If(
				!isDarwin,
				menu.RoledItem(menu.RoleAbout),
			),
		),
	)
}

var MenuBar, _ = markup.MustRegisterComponent((*menuBar)(nil))

type contextMenu struct {
	markup.Core
}

func (m *contextMenu) OnClickItem(e *html.MouseEvent) {
	log.PrintInfo("select Item: %#v", e.Target.ElementID)
	//TODO: Not implment menuitem.GetProp(), GetAttr(), etc...
}

func (m *contextMenu) Render() markup.RenderResult {
	return menu.ContextMenu(
		menu.Item("Item1", m.OnClickItem),
		menu.Item("Item2", m.OnClickItem),
	)
}

var (
	contextMenuInst *menu.MenuInstance
	ContextMenu, _  = markup.MustRegisterComponent((*contextMenu)(nil))
)

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
	err = contextMenuInst.Popup(e.ScreenPos(), w)
	if err != nil {
		panic(err)
	}
}

var TestChildComponent, _ = markup.MustRegisterComponent((*testChildComponent)(nil))

type testComponent struct {
	markup.Core
	Text    string `exciton:"text"`
	checked bool   `exciton:"checked"`
}

func (c *testComponent) clickHandler(e *html.MouseEvent) {
	log.PrintInfo("clickHandler is called!!")
	if c.Context().Builder() == nil {
		panic("Builder is NULL!!!!!!")
	}
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

var TestComponent, _ = markup.MustRegisterComponent((*testComponent)(nil))

func onAppStart() {
	log.PrintInfo("onAppStart")
	contextMenuInst = menu.MustNew(ContextMenu())
	menu.SetApplicationMenu(menu.MustNew(MenuBar()))

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
