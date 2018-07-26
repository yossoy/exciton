package mac

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"

	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/menu"
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
	callback(event.NewValueResult(event.NewJSONEncodedValueByEncodedBytes(jsonstr)))
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
	drvEvt := driver.DriverEvent{
		Name:      e.Name,
		Argument:  arg,
		Parameter: e.Params,
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

func (d *mac) relayEventWithResultToNative(e *event.Event, respCallback event.ResponceCallback) {
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
	if C.Driver_EmitEvent(C.CBytes(jb), C.NSUInteger(len(jb))) == 0 {
		driverLogError("Error: Driver_EmitEvent")
		//TODO: error responce!
	}
}

func (d *mac) requestEventEmit(devt *driver.DriverEvent) error {
	driverLogInfo("requestEventEmit: %v", devt)
	if devt.ResponceCallbackNo < 0 {
		v := event.NewJSONEncodedValueByEncodedBytes(devt.Argument)
		return event.Emit(devt.Name, v)
	}
	panic("not implement yet.")
}

func (d *mac) Init() error {
	C.Log_Init()

	driverLogInfo("driverMac.Init called")

	err := event.AddHandler("/app/quit", func(e *event.Event) {
		driverLogInfo("driver::terminate!!")
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

func (d *mac) IsIE() bool {
	return false
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
			{{Role: menu.RoleAbout}},
			{{Label: "services", Role: menu.RoleServices}},
			{
				{Role: menu.RoleHideOthers},
				{Role: menu.RoleUnhide},
			},
			{{Role: menu.RoleQuit}},
		}},
	{Label: "Edit",
		SubMenu: menu.MenuTemplate{
			{
				{Role: menu.RoleUndo},
				{Role: menu.RoleRedo},
			},
			{
				{Role: menu.RoleCut},
				{Role: menu.RoleCopy},
				{Role: menu.RolePaste},
				{Role: menu.RolePasteAndMatchStyle},
				{Role: menu.RoleDelete},
			},
			{
				{Role: menu.RoleStartSpeaking},
				{Role: menu.RoleStopSpeaking},
			},
		}},
	{Label: "Window", Role: menu.RoleWindow,
		SubMenu: menu.MenuTemplate{
			{
				{Role: menu.RoleMinimize},
				{Role: menu.RoleClose},
				{Role: menu.RoleFront},
			},
		}},
}

func Startup(startup app.StartupFunc) error {
	runtime.LockOSThread()
	event.StartEventMgr()
	defer event.StopEventMgr()
	si := &app.StartupInfo{
		AppMenu: macDefaultAppmenu,
	}
	d := newDriver()
	if err := d.Init(); err != nil {
		return err
	}
	return driver.Startup(d, &si.StartupInfo, func() error { return startup(si) })
}
