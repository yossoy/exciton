package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/yossoy/exciton/lang"

	"github.com/gorilla/websocket"

	"github.com/yossoy/exciton"
	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"

	// ievent "github.com/yossoy/exciton/internal/event"
	"github.com/yossoy/exciton/internal/markup"
	"github.com/yossoy/exciton/menu"
	"github.com/yossoy/exciton/window"
)

type sendMessageItem struct {
	Sync bool               `json:"sync"`
	Data driver.DriverEvent `json:"data"`
}

type web struct {
	quitChan        chan bool
	sendChan        chan []byte
	running         bool
	lock            *sync.Mutex
	respCallbacks   []event.ResponceCallback
	lastCallbackPos int
	serializer      *event.EventSerializer
}

type appOwner struct {
	preferredLanguages lang.PreferredLanguages
	ws                 *websocket.Conn
	sendChan           chan *sendMessageItem
}

func (ao *appOwner) PreferredLanguages() lang.PreferredLanguages {
	return ao.preferredLanguages
}

func (ao *appOwner) sendMessages() {
	for {
		d, ok := <-ao.sendChan
		if !ok {
			break
		}
		driverLogDebug("sendMessages: %v", d)
		ao.ws.WriteJSON(d)
	}
}

var (
	platform *web
)

func (d *web) addRespCallbackCallback(callback event.ResponceCallback) int {
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

func (d *web) responceCallback(jsonstr []byte, responceNo int) {
	d.lock.Lock()
	callback := d.respCallbacks[responceNo]
	d.respCallbacks[responceNo] = nil
	defer d.lock.Unlock()
	driverLogDebug("responceEventResult: %d => %v", responceNo, string(jsonstr))
	callback(event.NewValueResult(event.NewJSONEncodedValueByEncodedBytes(jsonstr)))
}

func (d *web) responceCallbackByValue(value interface{}, responceNo int) {
	d.lock.Lock()
	callback := d.respCallbacks[responceNo]
	d.respCallbacks[responceNo] = nil
	defer d.lock.Unlock()
	driverLogDebug("responceEventResult: %d => %v", responceNo, value)
	callback(event.NewValueResult(event.NewValue(value)))
}

func (d *web) responceCallbackError(err error, responceNo int) {
	d.lock.Lock()
	callback := d.respCallbacks[responceNo]
	d.respCallbacks[responceNo] = nil
	defer d.lock.Unlock()
	driverLogDebug("responceCallbackError: %d => %v", responceNo, err)
	callback(event.NewErrorResult(err))
}

func (d *web) relayEventToNative(e *event.Event) {
	a, err := app.GetAppFromEvent(e)
	if err != nil {
		panic(fmt.Errorf("app [%v] not found(%v)", e.Target, err))
	}
	ao := a.Owner().(*appOwner)
	var arg []byte
	if e.Argument != nil {
		arg, err = e.Argument.Encode()
		if err != nil {
			panic(err)
		}
	}
	// omit app Path?
	drvEvtPath, params := event.ToDriverEventPath(e.Target)
	driverLogDebug("relayEventToNative: drvEvtPath = %q, parms = %v", drvEvtPath, params)
	item := &sendMessageItem{
		Sync: false,
		Data: driver.DriverEvent{
			TargetPath: drvEvtPath,
			Name:       e.Name,
			Argument:   arg,
			Parameter:  params,
		},
	}
	ao.sendChan <- item
}

func (d *web) relayEventWithResultToNative(e *event.Event, respCallback event.ResponceCallback) {
	a, err := app.GetAppFromEvent(e)
	if err != nil {
		panic(fmt.Errorf("app [%v] not found (%v)", e.Target, err))
	}
	ao := a.Owner().(*appOwner)
	arg, err := e.Argument.Encode()
	if err != nil {
		panic(err)
	}
	drvEvtPath, params := event.ToDriverEventPath(e.Target)
	driverLogDebug("relayEventWithResultToNative: drvEvtPath = %q, parms = %v", drvEvtPath, params)

	item := &sendMessageItem{
		Sync: true,
		Data: driver.DriverEvent{
			TargetPath: drvEvtPath,
			Name:       e.Name,
			Argument:   arg,
			Parameter:  params,
			ResponceCallbackNo: d.addRespCallbackCallback(func(result event.Result) {
				driverLogDebug("responce...........%v\n", result)
				respCallback(result)
			}),
		},
	}
	ao.sendChan <- item
}

func (d *web) Init() error {
	app.AppClass.AddHandler("quit", func(e *event.Event) error {
		_, err := app.GetAppFromEvent(e)
		if err == nil {
			driverLogDebug("driver::terminate!!")
			//TODO: これ変だな。。。 platform.quitChanに送って良いのは全てのchannelが終了した時だけの筈
			platform.quitChan <- true
		}
		return nil
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

	err = initializeDialog(d.serializer)
	if err != nil {
		return err
	}
	return nil
}

func (d *web) Run() {
	d.running = true
	<-platform.quitChan
}

func createDirIfNotExists(name string) {
	_, err := os.Stat(name)
	if err != nil {
		os.MkdirAll(name, os.ModeDir|0755)
		return
	}
}

func (d *web) DriverType() string {
	return "web"
}

func (d *web) ResourcesFileSystem() (http.FileSystem, error) {
	resources, err := resourcesPath()
	if err != nil {
		return nil, err
	}
	return http.Dir(resources), nil
}

func (d *web) initEvent(si *app.StartupInfo) {
	app.InitEvents(false, si)

	d.serializer = event.NewSerializer(d.relayEventToNative, d.relayEventWithResultToNative)
}

func resourcesPath() (string, error) {
	exePathStr, err := os.Executable()
	if err != nil {
		return "", err
	}
	resourcesName := filepath.Join(filepath.Dir(exePathStr), "resources")
	//TODO: ??? need to create folder?
	createDirIfNotExists(resourcesName)
	return resourcesName, nil
}

func (d *web) NativeRequestJSMethod() string {
	return "window.parent.exciton.callWindowMethod"
}

func (d *web) Log(lvl driver.LogLevel, msg string, args ...interface{}) {
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

func newDriver() *web {
	platform = &web{
		quitChan: make(chan bool),
		lock:     new(sync.Mutex),
	}
	return platform
}

func internalInitFunc(app *app.App, info *app.StartupInfo) error {
	driverLogDebug("initializeInitFunc: 1")
	menu.SetApplicationMenu(app, info.AppMenu)
	driverLogDebug("initializeInitFunc: 2")
	if info.OnAppStart != nil {
		err := info.OnAppStart(app, info)
		if err != nil {
			return err
		}
	}
	driverLogDebug("initializeInitFunc: 3")
	cfg := &window.WindowConfig{}
	var rr markup.RenderResult
	if info.OnNewWindow != nil {
		var err error
		rr, err = info.OnNewWindow(app, cfg)
		if err != nil {
			return err
		}
	} else {
		rr = emptyPage()
	}
	win, err := window.NewWindow(app, cfg)
	if err != nil {
		return err
	}
	app.MainWindow = win
	win.Mount(rr)
	return nil
}

func emptyPage() markup.RenderResult {
	return html.Div(
		markup.Text("Empty"),
	)
}

func webRoot(id string) ([]byte, error) {
	ctx := struct {
		ID string
	}{
		ID: id,
	}
	f, err := fileSystem.Open("/webroot.gohtml")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	t, err := template.New("").Parse(string(b))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, ctx)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func rootHTMLHandler(si *app.StartupInfo) http.HandlerFunc {
	// redirect or iframe?
	return func(w http.ResponseWriter, r *http.Request) {
		driverLogDebug("rootHTMLHandler: %q", r.RequestURI)
		ao := &appOwner{}
		if al := r.Header.Get("Accept-Language"); al != "" {
			pl, err := lang.NewPreferredLanguagesFromAcceptLanguages(al)
			if err != nil {
				driverLogDebug("webRoot faild: %v", err)
				http.Error(w, err.Error(), http.StatusProcessing)
				return
			}
			ao.preferredLanguages = pl
		}
		a := app.NewApp(ao)
		pb, err := webRoot(a.TargetID())
		if err != nil {
			driverLogDebug("webRoot faild: %v", err)
			http.Error(w, err.Error(), http.StatusProcessing)
			return
		}
		driverLogDebug("webroot => %s", string(pb))
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(pb)
	}
}

func driverWebSockHandler(si *app.StartupInfo) http.HandlerFunc {
	upgrader := websocket.Upgrader{}
	return func(w http.ResponseWriter, r *http.Request) {
		driverLogDebug("driverWebSockHandler!!")
		appid, ok := driver.RequestVars(r)["appid"]
		if !ok {
			driverLogDebug("invalid path")
			http.Error(w, "invalid path", http.StatusNotFound)
			return
		}
		a := app.GetAppByID(appid)
		if a == nil {
			driverLogDebug("invalid appid: %v", appid)
			http.Error(w, "invalid appid", http.StatusProcessing)
			return
		}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			driverLogError("upgrade error: %v", err)
			http.Error(w, err.Error(), http.StatusProcessing)
			return
		}
		ao := a.Owner().(*appOwner)
		ao.ws = c
		ao.sendChan = make(chan *sendMessageItem)
		defer c.Close()
		go ao.sendMessages()
		for {
			var devt driver.DriverEvent
			if err := c.ReadJSON(&devt); err != nil {
				driverLogError("driverWebSocket: read error: %v", err)
				break
			}
			driverLogDebug("driverEvent ==> %v", devt)
			if devt.ResponceCallbackNo < 0 {
				v := event.NewJSONEncodedValueByEncodedBytes(devt.Argument)
				et, err := event.StringToEventTarget(devt.TargetPath)
				if err != nil {
					driverLogDebug("event.Emit drv event error: %v", err)
					panic(err)
				}
				go func() {
					driverLogDebug("target: %v, name: %q, arg: %v", et, devt.Name, v)
					if err = event.Emit(et, devt.Name, v); err != nil {
						driverLogDebug("event.Emit error: %v", err)
						//panic(err)
					}
				}()
			} else {
				if devt.TargetPath == "" && devt.Name == "responceEventResult" {
					//TODO: error
					driverLogDebug("/responceEventResult: %s", string(devt.Argument))
					r := event.NewValueResult(event.NewJSONEncodedValueByEncodedBytes(devt.Argument))
					callback := platform.respCallbacks[devt.ResponceCallbackNo]
					platform.respCallbacks[devt.ResponceCallbackNo] = nil
					go callback(r)
				} else {
					panic(fmt.Errorf("not implement yet"))
				}
			}
		}
	}
}

func addWebRootResources(si *driver.StartupInfo) error {
	si.Router.PathPrefix("/exciton/web/assets/").Handler(http.StripPrefix("/exciton/web/assets/", http.FileServer(fileSystem)))
	return nil
}

// Startup is startup function in windows.
func Startup(startup app.StartupFunc) error {
	runtime.LockOSThread()
	//	ievent.StartEventMgr()
	//	defer ievent.StopEventMgr()
	si := &app.StartupInfo{}
	si.StartupInfo.PortNo = 8080
	si.StartupInfo.AppURLBase = "/exciton/{appid}"
	d := newDriver()
	d.initEvent(si)
	defer d.serializer.Stop()

	if err := d.Init(); err != nil {
		return err
	}
	sf := func() error {
		driverLogDebug("called init func: 1")
		if err := startup(si); err != nil {
			return err
		}
		driverLogDebug("called init func: 2")
		if err := exciton.Init(si, internalInitFunc); err != nil {
			return err
		}
		driverLogDebug("called init func: 3")
		si.Router.HandleFunc("/app/{appid}/ws", driverWebSockHandler(si))
		si.Router.HandleFunc("/", rootHTMLHandler(si))
		driverLogDebug("called init func: 4")
		return nil
	}
	driver.AddInitProc(driver.InitProcTimingPostStartup, addWebRootResources)
	return driver.Startup(d, &si.StartupInfo, sf)
}
