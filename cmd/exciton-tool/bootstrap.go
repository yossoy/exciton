package main

const mainTempl = `package main

import (
	"github.com/yossoy/exciton"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
)

func onAppStart(app *exciton.App) {
	menu.SetApplicationMenu(menu.MustNew(MenuBar()))

	cfg := window.WindowConfig{
		Title: "Sample",
	}
	w, err := window.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	w.Mount(MainView())
}

func main() {
	exciton.Run(onAppStart)
}
`

const menuTempl = `package main

import (
	"runtime"

	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/menu"
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

// MenuBar is menubar component.
var MenuBar = markup.MustRegisterComponent((*menuBar)(nil))
`

const viewTempl = `package main

import (
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/markup"
)

type mainView struct {
	markup.Core
}

func (mv *mainView) Render() markup.RenderResult {
	return html.Div(
		markup.Text("Sample"),
	)
}

// MainView is main view component.
var MainView = markup.MustRegisterComponent((*mainView)(nil))
`

const driverMainTempl = `package main

import (
	_ "github.com/yossoy/exciton/driver/{{.DriverName}}"
)
`
