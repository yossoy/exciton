package mac

import (
	"strings"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"

	"github.com/yossoy/exciton"
	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/internal/markup"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
	"github.com/yossoy/exciton/lang"
)

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa
#cgo LDFLAGS: -framework WebKit
#cgo LDFLAGS: -framework CoreImage
#cgo LDFLAGS: -framework Security

#include "driver.h"
#include "log.h"
*/
import "C"

type mac struct {
	running         bool
	lock            *sync.Mutex
	respCallbacks   []event.ResponceCallback
	lastCallbackPos int
	serializer *event.EventSerializer
}

var (
	platform *mac
)

func (d *mac) addRespCallbackCallback(callback event.ResponceCallback) int {
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

func (d *mac) responceCallback(jsonstr []byte, responceNo int) {
	d.lock.Lock()
	callback := d.respCallbacks[responceNo]
	d.respCallbacks[responceNo] = nil
	defer d.lock.Unlock()
	driverLogDebug("responceEventResult: %d => %v", responceNo, string(jsonstr))
	go callback(event.NewValueResult(event.NewJSONEncodedValueByEncodedBytes(jsonstr)))
}

func (d *mac) relayEventToNative(e *event.Event) {
	var arg []byte
	var err error
	if e.Argument != nil {
		arg, err = e.Argument.Encode()
		if err != nil {
			panic(err)
		}
	}
	drvEvtPath, params := event.ToDriverEventPath(e.Target, e.Name)
	driverLogDebug("relayEventToNative: drvEvtPath = %q, parms = %v", drvEvtPath, params)
	drvEvt := driver.DriverEvent{
		Name:      drvEvtPath,
		Argument:  arg,
		Parameter: params,
	}
	jb, err := json.Marshal(&drvEvt)
	if err != nil {
		panic(err)
	}
	driverLogDebug("emitToNative: %s", string(jb))
	if C.Driver_EmitEvent(C.CBytes(jb), C.NSUInteger(len(jb))) == 0 {
		driverLogError("Error: Driver_EmitEvent")
		//TODO: error responce!
	}
}

func (d *mac) relayEventWithResultToNative(e *event.Event, respCallback event.ResponceCallback) {
	arg, err := e.Argument.Encode()
	if err != nil {
		panic(err)
	}
	drvEvtPath, params := event.ToDriverEventPath(e.Target, e.Name)
	driverLogDebug("replayEventToNative: drvEvtPath = %q, parms = %v", drvEvtPath, params)

	drvEvt := driver.DriverEvent{
		Name:      drvEvtPath,
		Argument:  arg,
		Parameter: params,
		ResponceCallbackNo: d.addRespCallbackCallback(func(result event.Result) {
			driverLogDebug("responce...........%v\n", result)
			respCallback(result)
		}),
	}
	jb, err := json.Marshal(&drvEvt)
	if err != nil {
		panic(err)
	}
	if C.Driver_EmitEvent(C.CBytes(jb), C.NSUInteger(len(jb))) == 0 {
		driverLogError("Error: Driver_EmitEvent")
		//TODO: error responce!
	}
}

func (d *mac) requestEventEmit(devt *driver.DriverEvent) error {
	driverLogInfo("requestEventEmit: %v", devt)
	if devt.ResponceCallbackNo < 0 {
		v := event.NewJSONEncodedValueByEncodedBytes(devt.Argument)
		t, name, err := event.StringToEventTarget(devt.Name)
		if err != nil {
			return err
		}
		go event.Emit(t, name, v)
		return nil
	}
	panic("not implement yet.")
}

func (d *mac) Init() error {
	C.Log_Init()

	driverLogInfo("driverMac.Init called")

	app.AppClass.AddHandler("quit", func(e *event.Event) {
		driverLogInfo("driver::terminate!!")
		C.Driver_Terminate()
	})

	var err error

	err = initializeWindow(d.serializer)
	if err != nil {
		return err
	}

	err = initializeMenu(d.serializer)
	if err != nil {
		return err
	}

	initializeDialog(d.serializer)
	return nil
}

func (d *mac) Run() {
	d.running = true
	C.Driver_Run()
}

func isAppBundled(exePathStr string) (bool, error) {
	dir, err := filepath.Abs(filepath.Dir(exePathStr))
	if err != nil {
		return false, err
	}
	for {
		baseName := filepath.Base(dir)
		if baseName == dir {
			break
		}
		if filepath.Ext(baseName) == ".app" {
			return true, nil
		}
		dir = filepath.Dir(dir)
	}
	return false, nil
}

func createDirIfNotExists(name string) {
	_, err := os.Stat(name)
	if err != nil {
		os.MkdirAll(name, os.ModeDir|0755)
		return
	}
}

func (d *mac) DriverType() string {
	return "mac"
}

func (d *mac) ResourcesFileSystem() (http.FileSystem, error) {
	resources, err := resourcesPath()
	if err != nil {
		return nil, err
	}
	return http.Dir(resources), nil
}

func resourcesPath() (string, error) {
	exePathStr, err := os.Executable()
	if err != nil {
		return "", err
	}
	bundled, err := isAppBundled(exePathStr)
	if err != nil {
		return "", err
	}
	if bundled {
		cresources := C.Driver_GetBundleResourcesPath()
		resourcesName := C.GoString(cresources)
		C.free(unsafe.Pointer(cresources))
		return filepath.Join(resourcesName, "assets"), nil

	}
	resourcesName := filepath.Join(filepath.Dir(exePathStr), "resources")
	//TODO: ??? need to create folder?
	createDirIfNotExists(resourcesName)
	return resourcesName, nil
}

func preferredLanguages() lang.PreferredLanguages {
	clangs := C.Driver_GetPreferrdLanguage();
	langs := C.GoString(clangs);
	C.free(unsafe.Pointer(clangs));
	langAndLocales := strings.Split(langs, ";")
	return lang.NewPreferredLanguages(langAndLocales...)
}

func (d *mac) NativeRequestJSMethod() string {
	return "webkit.messageHandlers.golangRequest.postMessage"
}

func (d *mac) Log(lvl driver.LogLevel, msg string, args ...interface{}) {
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

func newDriver() *mac {
	platform = &mac{
		lock: new(sync.Mutex),
	}
	return platform
}

//export requestEventEmit
func requestEventEmit(cstr unsafe.Pointer, clen C.int) {
	jsonstr := C.GoBytes(cstr, clen)
	devt := driver.DriverEvent{}
	driverLogDebug("requestEventEmit: %s", jsonstr)
	if err := json.Unmarshal(jsonstr, &devt); err != nil {
		panic(err)
	}
	err := platform.requestEventEmit(&devt)
	if err != nil {
		if devt.ResponceCallbackNo >= 0 {
			// send error responce
			panic(err)
		}
		driverLogError("event emit failed: %q : %v", err, devt)
		//TODO: error?
	}
}

//export responceEventResult
func responceEventResult(crespNo C.int, cstr unsafe.Pointer, clen C.int) {
	jsonstr := C.GoBytes(cstr, clen)
	respNo := int(crespNo)

	platform.responceCallback(jsonstr, respNo)
}

var macDefaultAppmenu = menu.AppMenuTemplate{
	{Label: menu.AppMenuLabel,
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
	{Label: "Edit",
		SubMenu: menu.MenuTemplate{
			{Role: menu.RoleUndo},
			{Role: menu.RoleRedo},
			menu.Separator,
			{Role: menu.RoleCut},
			{Role: menu.RoleCopy},
			{Role: menu.RolePaste},
			{Role: menu.RolePasteAndMatchStyle},
			{Role: menu.RoleDelete},
			menu.Separator,
			{Role: menu.RoleStartSpeaking},
			{Role: menu.RoleStopSpeaking},
		}},
	{Label: "Window", Role: menu.RoleWindow,
		SubMenu: menu.MenuTemplate{
			{Role: menu.RoleMinimize},
			{Role: menu.RoleClose},
			{Role: menu.RoleFront},
		}},
}

func emptyPage() markup.RenderResult {
	return html.Div(
		markup.Text("Empty"),
	)
}

func internalInitFunc(a *app.App, info *app.StartupInfo) error {
	driverLogDebug("internalInitFunc1: %v", a)
	menu.SetApplicationMenu(a, info.AppMenu)
	if info.OnAppStart != nil {
		err := info.OnAppStart(a, info)
		if err != nil {
			return err
		}
	}
	driverLogDebug("internalInitFunc2: %v", a)

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
	driverLogDebug("internalInitFunc3: %v", a)
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

func (d *mac) initEvent(si *app.StartupInfo) {
	app.InitEvents(true, si)

	d.serializer = event.NewSerializer(d.relayEventToNative, d.relayEventWithResultToNative)
}

func Startup(startup app.StartupFunc) error {
	runtime.LockOSThread()
	si := &app.StartupInfo{
		AppMenu: macDefaultAppmenu,
	}
	d := newDriver()
	d.initEvent(si)
	defer d.serializer.Stop()
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
