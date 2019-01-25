package windows

import (
	"github.com/yossoy/exciton/lang"
	"strings"
	"encoding/json"
	"os"
	"runtime"
	"sync"
	"unsafe"

	"github.com/yossoy/exciton"
	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"
	ievent "github.com/yossoy/exciton/internal/event"
	"github.com/yossoy/exciton/markup"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
)

/*
#include <stdlib.h>
#include "driver.h"
#include "log.h"
*/
import "C"

type windows struct {
	running         bool
	lock            *sync.Mutex
	respCallbacks   []event.ResponceCallback
	lastCallbackPos int
}

var (
	platform *windows
)

func (d *windows) addRespCallbackCallback(callback event.ResponceCallback) int {
	d.lock.Lock()
	defer d.lock.Unlock()
	for i := 0; i < len(d.respCallbacks); i++ {
		idx := (d.lastCallbackPos + i) % len(d.respCallbacks)
		if d.respCallbacks[idx] == nil {
			d.respCallbacks[idx] = callback
			d.lastCallbackPos = idx
			return idx
		}
	}
	idx := len(d.respCallbacks)
	d.respCallbacks = append(d.respCallbacks, callback)
	d.lastCallbackPos = 0
	return idx
}

func (d *windows) responceCallback(jsonstr []byte, responceNo int) {
	d.lock.Lock()
	callback := d.respCallbacks[responceNo]
	d.respCallbacks[responceNo] = nil
	defer d.lock.Unlock()
	driverLogDebug("responceEventResult: %d => %v", responceNo, string(jsonstr))
	callback(event.NewValueResult(event.NewJSONEncodedValueByEncodedBytes(jsonstr)))
}

func (d *windows) relayEventToNative(e *event.Event) {
	var arg []byte
	var err error
	if e.Argument != nil {
		arg, err = e.Argument.Encode()
		if err != nil {
			panic(err)
		}
	}
	drvEvt := driver.DriverEvent{
		Name:      e.Name,
		Argument:  arg,
		Parameter: e.Params,
	}
	jb, err := json.Marshal(&drvEvt)
	if err != nil {
		panic(err)
	}
	driverLogDebug("relayEventToNative: %q", jb)
	if C.Driver_EmitEvent(C.CBytes(jb), C.int(len(jb))) == 0 {
		driverLogError("Error: Driver_EmitEvent")
	}
}

func (d *windows) relayEventWithResultToNative(e *event.Event, respCallback event.ResponceCallback) {
	arg, err := e.Argument.Encode()
	if err != nil {
		panic(err)
	}
	drvEvt := driver.DriverEvent{
		Name:      e.Name,
		Argument:  arg,
		Parameter: e.Params,
		ResponceCallbackNo: d.addRespCallbackCallback(func(result event.Result) {
			driverLogDebug("responce...........%v\n", result)
			respCallback(result)
		}),
	}
	jb, err := json.Marshal(&drvEvt)
	if err != nil {
		panic(err)
	}
	driverLogDebug(string(jb))
	if C.Driver_EmitEvent(C.CBytes(jb), C.int(len(jb))) == 0 {
		driverLogError("Error: Driver_EmitEvent")
	}
}

func (d *windows) requestEventEmit(devt *driver.DriverEvent) error {
	driverLogDebug("requestEventEmit: %v", devt)
	if devt.ResponceCallbackNo < 0 {
		v := event.NewJSONEncodedValueByEncodedBytes(devt.Argument)
		return ievent.Emit(devt.Name, v)
	}
	panic("not implement yet.")
}

func (d *windows) Init() error {
	C.Log_Init()

	err := ievent.AddHandler("/app/quit", func(e *event.Event) {
		driverLogDebug("driver::terminate!!")
		C.Driver_Terminate()
	})

	err = initializeWindow()
	if err != nil {
		return err
	}

	err = initializeMenu()
	if err != nil {
		return err
	}

	err = initializeDialog()
	if err != nil {
		return err
	}
	return nil
}

func (d *windows) Run() {
	d.running = true
	//TODO: emit /init in native code
	ievent.Emit("/app/init", event.NewValue(nil))
	C.Driver_Run()
}

func createDirIfNotExists(name string) {
	_, err := os.Stat(name)
	if err != nil {
		os.MkdirAll(name, os.ModeDir|0755)
		return
	}
}

func (d *windows) DriverType() string {
	return "ie"
}

func preferredLanguages() lang.PreferredLanguages {
	clangs := C.Driver_GetPreferrdLanguage();
	langs := C.GoString(clangs);
	C.free(unsafe.Pointer(clangs));
	langAndLocales := strings.Split(langs, ";")
	return lang.NewPreferredLanguages(langAndLocales...)
}

func (d *windows) NativeRequestJSMethod() string {
	return "window.external.golangRequest"
}

func (d *windows) Log(lvl driver.LogLevel, msg string, args ...interface{}) {
	switch lvl {
	case driver.LogLevelDebug:
		driverLogDebug(msg, args...)
	case driver.LogLevelInfo:
		driverLogInfo(msg, args...)
	case driver.LogLevelWarning:
		driverLogWarning(msg, args...)
	case driver.LogLevelError:
		driverLogError(msg, args...)
	}
}

func newDriver() *windows {
	platform = &windows{
		lock: new(sync.Mutex),
	}
	return platform
}

//export requestEventEmit
func requestEventEmit(cstr unsafe.Pointer, clen C.int) {
	jsonstr := C.GoBytes(cstr, clen)
	devt := driver.DriverEvent{}
	if err := json.Unmarshal(jsonstr, &devt); err != nil {
		panic(err)
	}
	err := platform.requestEventEmit(&devt)
	if err != nil {
		if devt.ResponceCallbackNo >= 0 {
			// send error responce
			panic(err)
		}
		driverLogDebug("event emit failed: %q : %v", err, devt)
	}
}

//export responceEventResult
func responceEventResult(crespNo C.int, cstr unsafe.Pointer, clen C.int) {
	jsonstr := C.GoBytes(cstr, clen)
	respNo := int(crespNo)

	platform.responceCallback(jsonstr, respNo)
}

var windowsDefaultMenu = menu.AppMenuTemplate{
	{Label: "File",
		SubMenu: menu.MenuTemplate{
			{
				{Role: menu.RoleClose},
			},
			{
				{Role: menu.RoleQuit},
			},
		}},
	{Label: "Edit",
		SubMenu: menu.MenuTemplate{
			{
				{Role: menu.RoleCut},
				{Role: menu.RoleCopy},
				{Role: menu.RolePaste},
				{Role: menu.RoleDelete},
			},
		}},
	{Label: "Help", Role: menu.RoleHelp,
		SubMenu: menu.MenuTemplate{
			{{Role: menu.RoleAbout}},
		}},
}

func emptyPage() markup.RenderResult {
	return html.Div(
		markup.Text("Empty"),
	)
}

func internalInitFunc(a *app.App, info *app.StartupInfo) error {
	menu.SetApplicationMenu(a, info.AppMenu)
	if info.OnAppStart != nil {
		err := info.OnAppStart(a, info)
		if err != nil {
			return err
		}
	}

	winCfg := &window.WindowConfig{}
	var rr markup.RenderResult
	if info.OnNewWindow != nil {
		var err error
		rr, err = info.OnNewWindow(a, winCfg)
		if err != nil {
			return err
		}
	} else {
		rr = emptyPage()
	}
	w, err := window.NewWindow(a, winCfg)
	if err != nil {
		return err
	}
	a.MainWindow = w
	return w.Mount(rr)
}

type appOwner struct {
	preferredLanguages lang.PreferredLanguages
}

func (ao *appOwner) PreferredLanguages() lang.PreferredLanguages {
	return ao.preferredLanguages
}

// Startup is startup function in windows.
func Startup(startup app.StartupFunc) error {
	runtime.LockOSThread()
	ievent.StartEventMgr()
	defer ievent.StopEventMgr()
	si := &app.StartupInfo{
		AppMenu: windowsDefaultMenu,
	}
	si.StartupInfo.AppEventRoot = ievent.RootGroup()
	d := newDriver()
	if err := d.Init(); err != nil {
		return err
	}
	sf := func() error {
		if err := startup(si); err != nil {
			return err
		}
		appOwner := &appOwner {
			preferredLanguages: preferredLanguages(),
		}
		app.NewSingletonApp(appOwner)
		if err := exciton.Init(si, internalInitFunc); err != nil {
			return err
		}
		return nil
	}
	return driver.Startup(d, &si.StartupInfo, sf)
}
