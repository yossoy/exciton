package windows

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"

	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
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
	if C.Driver_EmitEvent(C.CBytes(jb), C.int(len(jb))) == 0 {
		driverLogError("Error: Driver_EmitEvent")
	}
}

func (d *windows) requestEventEmit(devt *driver.DriverEvent) error {
	driverLogDebug("requestEventEmit: %v", devt)
	if devt.ResponceCallbackNo < 0 {
		v := event.NewJSONEncodedValueByEncodedBytes(devt.Argument)
		return event.Emit(devt.Name, v)
	}
	panic("not implement yet.")
}

func (d *windows) Init() error {
	C.Log_Init()

	err := event.AddHandler("/app/quit", func(e *event.Event) {
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
	event.Emit("/app/init", event.NewValue(nil))
	C.Driver_Run()
}

func createDirIfNotExists(name string) {
	_, err := os.Stat(name)
	if err != nil {
		os.MkdirAll(name, os.ModeDir|0755)
		return
	}
}

func (d *windows) IsIE() bool {
	return true
}

func (d *windows) Resources() (string, error) {
	exePathStr, err := os.Executable()
	if err != nil {
		return "", err
	}
	resourcesName := filepath.Join(filepath.Dir(exePathStr), "resources")
	createDirIfNotExists(resourcesName)
	return resourcesName, nil
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

func init() {
	runtime.LockOSThread()
	driver.SetupDriver(newDriver())
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
