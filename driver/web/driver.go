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

	"github.com/gorilla/websocket"

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
}

type webAppDriverData struct {
	ws       *websocket.Conn
	sendChan chan *sendMessageItem
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
	appid, ok := e.Params["appid"]
	if !ok {
		panic("invalid event (appid not found)")
	}
	a := app.GetAppByID(appid)
	if a == nil {
		panic(fmt.Errorf("app [%s] not found", appid))
	}
	dd := a.DriverData.(*webAppDriverData)
	var arg []byte
	var err error
	if e.Argument != nil {
		arg, err = e.Argument.Encode()
		if err != nil {
			panic(err)
		}
	}
	item := &sendMessageItem{
		Sync: false,
		Data: driver.DriverEvent{
			Name:      e.Name,
			Argument:  arg,
			Parameter: e.Params,
		},
	}
	dd.sendChan <- item
}

func (d *web) relayEventWithResultToNative(e *event.Event, respCallback event.ResponceCallback) {
	appid, ok := e.Params["appid"]
	if !ok {
		panic("invalid event (appid not found)")
	}
	a := app.GetAppByID(appid)
	if a == nil {
		panic(fmt.Errorf("app [%s] not found", appid))
	}
	dd := a.DriverData.(*webAppDriverData)
	arg, err := e.Argument.Encode()
	if err != nil {
		panic(err)
	}
	item := &sendMessageItem{
		Sync: true,
		Data: driver.DriverEvent{
			Name:      e.Name,
			Argument:  arg,
			Parameter: e.Params,
			ResponceCallbackNo: d.addRespCallbackCallback(func(result event.Result) {
				driverLogDebug("responce...........%v\n", result)
				respCallback(result)
			}),
		},
	}
	dd.sendChan <- item
}

func (d *web) Init() error {
	g, err := ievent.AddGroup("/exciton/:appid")
	if err != nil {
		return err
	}
	err = g.AddHandler("/app/quit", func(e *event.Event) {
		appid := e.Params["appid"]
		a := app.GetAppByID(appid)
		if a != nil {
			driverLogDebug("driver::terminate!!")
			platform.quitChan <- true
		}
	})
	if err != nil {
		return err
	}

	err = initializeWindow(g)
	if err != nil {
		return err
	}

	err = initializeMenu(g)
	if err != nil {
		return err
	}

	err = initializeDialog(g)
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
	menu.SetApplicationMenu("/exciton/"+app.ID, info.AppMenu)
	if info.OnAppStart != nil {
		err := info.OnAppStart(app, info)
		if err != nil {
			return err
		}
	}
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
	win, err := window.NewWindow(app.ID, cfg)
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
		dd := &webAppDriverData{}
		a := app.NewApp(dd)
		pb, err := webRoot(a.ID)
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

func (dd *webAppDriverData) sendMessages() {
	for {
		d, ok := <-dd.sendChan
		if !ok {
			break
		}
		driverLogDebug("sendMessages: %v", d)
		dd.ws.WriteJSON(d)
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
		dd := a.DriverData.(*webAppDriverData)
		dd.ws = c
		dd.sendChan = make(chan *sendMessageItem)
		defer c.Close()
		go dd.sendMessages()
		for {
			var devt driver.DriverEvent
			if err := c.ReadJSON(&devt); err != nil {
				driverLogError("driverWebSocket: read error: %v", err)
				break
			}
			driverLogDebug("driverEvent ==> %v", devt)
			if devt.ResponceCallbackNo < 0 {
				v := event.NewJSONEncodedValueByEncodedBytes(devt.Argument)
				if err := ievent.Emit(devt.Name, v); err != nil {
					driverLogDebug("event.Emit error: %v", err)
					//panic(err)
				}
			} else {
				if devt.Name == "/responceEventResult" {
					//TODO: error
					driverLogDebug("/responceEventResult: %s", string(devt.Argument))
					r := event.NewValueResult(event.NewJSONEncodedValueByEncodedBytes(devt.Argument))
					callback := platform.respCallbacks[devt.ResponceCallbackNo]
					platform.respCallbacks[devt.ResponceCallbackNo] = nil
					callback(r)
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
	ievent.StartEventMgr()
	defer ievent.StopEventMgr()
	si := &app.StartupInfo{}
	si.StartupInfo.PortNo = 8080
	si.StartupInfo.AppURLBase = "/exciton/{appid}"
	rootGroup, err := ievent.AddGroup("/exciton/:appid")
	if err != nil {
		return err
	}
	si.StartupInfo.AppEventRoot = rootGroup

	d := newDriver()
	if err := d.Init(); err != nil {
		return err
	}
	sf := func() error {
		if err := startup(si); err != nil {
			return err
		}
		if err := exciton.Init(si, internalInitFunc); err != nil {
			return err
		}
		si.Router.HandleFunc("/exciton/{appid}/ws", driverWebSockHandler(si))
		si.Router.HandleFunc("/", rootHTMLHandler(si))
		return nil
	}
	driver.AddInitProc(driver.InitProcTimingPostStartup, addWebRootResources)
	return driver.Startup(d, &si.StartupInfo, sf)
}
