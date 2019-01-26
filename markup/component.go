package markup

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/yossoy/exciton/driver"

	"github.com/yossoy/exciton/internal/markup"
)

type Core struct {
	markup.Core
}

type Component = markup.Component

type ComponentRegisterParameter = markup.ComponentRegisterParameter
type ComponentInstance = markup.ComponentInstance
type InitInfo = markup.InitInfo
type ClassInitProc = markup.ClassInitProc

func RegisterComponent(c Component, params ...ComponentRegisterParameter) (ComponentInstance, error) {
	_, fp, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("invalid caller")
	}
	ci, _, err := markup.RegisterComponent(c, filepath.Dir(fp), params, false)
	return ci, err
}

func MustRegisterComponent(c Component, params ...ComponentRegisterParameter) ComponentInstance {
	_, fp, _, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Errorf("invalid caller"))
	}
	ci, _, err := markup.RegisterComponent(c, filepath.Dir(fp), params, false)
	if err != nil {
		panic(err)
	}
	return ci
}

func RegisterKlassOnly(c Component, params ...ComponentRegisterParameter) (Klass, error) {
	_, fp, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("invalid caller")
	}
	_, k, err := markup.RegisterComponent(c, filepath.Dir(fp), params, true)
	return k, err
}

func MustRegisterKlassOnly(c Component, params ...ComponentRegisterParameter) Klass {
	_, fp, _, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Errorf("invalid caller"))
	}
	_, k, err := markup.RegisterComponent(c, filepath.Dir(fp), params, true)
	if err != nil {
		panic(err)
	}
	return k
}

func WithClassInitializer(timing driver.InitProcTiming, proc ClassInitProc) ComponentRegisterParameter {
	return markup.WithClassInitializer(timing, proc)
}

func WithGlobalStyleSheet(css string) ComponentRegisterParameter {
	return markup.WithGlobalStyleSheet(css)
}

func WithComponentStyleSheet(css string) ComponentRegisterParameter {
	return markup.WithComponentStyleSheet(css)
}

func WithGlobalScript(js string) ComponentRegisterParameter {
	return markup.WithGlobalScript(js)
}

func WithComponentScript(js string) ComponentRegisterParameter {
	return markup.WithComponentScript(js)
}
