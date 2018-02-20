package driver

import (
	"encoding/json"
	"errors"
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
	Init() error
	Run()
	IsIE() bool
	Resources() (string, error)
	NativeRequestJSMethod() string
	Log(lvl LogLevel, msg string, args ...interface{})
}

var (
	platform Driver
)

func SetupDriver(driver Driver) {
	platform = driver
}

func Init() error {
	if platform == nil {
		return errors.New("driver is not loaded.")
	}
	return platform.Init()
}

func Run() {
	platform.Run()
}

func IsIE() bool {
	return platform.IsIE()
}

func Log(lvl LogLevel, fmt string, args ...interface{}) {
	platform.Log(lvl, fmt, args...)
}

func Resources() (string, error) {
	return platform.Resources()
}

func NativeRequestJSMethod() string {
	return platform.NativeRequestJSMethod()
}
