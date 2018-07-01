package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

type DriverEvent struct {
	Name               string            `json:"name"`
	Argument           json.RawMessage   `json:"argument"`
	Parameter          map[string]string `json:"parameter"`
	ResponceCallbackNo int               `json:"respCallbackNo"`
}

type Driver interface {
	Run()
	IsIE() bool
	ResourcesFileSystem() (http.FileSystem, error)
	NativeRequestJSMethod() string
	Log(lvl LogLevel, msg string, args ...interface{})
}

var (
	platform Driver
)

func Run() {
	platform.Run()
}

func IsIE() bool {
	return platform.IsIE()
}

func Log(lvl LogLevel, fmt string, args ...interface{}) {
	platform.Log(lvl, fmt, args...)
}

func NativeRequestJSMethod() string {
	return platform.NativeRequestJSMethod()
}

type StartupInfo struct {
	PortNo     int
	Router     *mux.Router
	OnAppStart func()
	OnAppQuit  func()
}

type StartupFunc func(info *StartupInfo) error

func newStartupInfo() *StartupInfo {
	return &StartupInfo{
		Router: mux.NewRouter(),
	}
}

var BaseURL string

func ResourcesFileSystem() (http.FileSystem, error) {
	return platform.ResourcesFileSystem()
}

func Startup(driver Driver, startup StartupFunc) error {
	platform = driver
	si := newStartupInfo()
	if err := startup(si); err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", si.PortNo))
	if err != nil {
		return err
	}
	si.PortNo = l.Addr().(*net.TCPAddr).Port

	BaseURL = fmt.Sprintf(fmt.Sprintf("http://127.0.0.1:%d", si.PortNo))

	Log(LogLevelDebug, "**** Internal Web server port: %d", si.PortNo)

	fs, err := ResourcesFileSystem()
	if err != nil {
		return err
	}

	// logging middleware is inject in debug mode only.
	if !ReleaseBuild {
		si.Router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Log(LogLevelInfo, "Request[%s]:: %q", r.Method, r.RequestURI)
				next.ServeHTTP(w, r)
			})
		})
	}

	// add "/resources/" file handler
	si.Router.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", http.FileServer(fs)))

	srv := &http.Server{
		Handler: si.Router,
	}
	go func() {
		defer l.Close()
		srv.Serve(l)
		//http.Serve(l, si.Router)
	}()

	platform.Run()

	//TODO: MacでCommand+Qで終了した時にここに辿り着かない?

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		cancelFunc()
		return err
	}
	cancelFunc()

	return nil
}
