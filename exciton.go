package exciton

import "github.com/yossoy/exciton/event"
import "github.com/yossoy/exciton/driver"
import "github.com/yossoy/exciton/markup"
import "github.com/yossoy/exciton/menu"
import "github.com/yossoy/exciton/window"

// RunCallback is called at ready application
type RunCallback func(app *App)

// Run start application mainloop
func Run(callback RunCallback) {
	event.StartEventMgr()
	event.AddHandler("/app/init", func(e *event.Event) {
		callback(&app)
	})
	err := driver.Init()
	if err != nil {
		panic(err)
	}
	err = window.InitWindows()
	if err != nil {
		panic(err)
	}
	err = menu.InitMenus()
	if err != nil {
		panic(err)
	}
	err = markup.InitEvents()
	if err != nil {
		panic(err)
	}
	driver.Run()
}
