package markup

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	componentsLock   = sync.RWMutex{}
	componentKlasses = make(map[string]*Klass)
)

type Klass struct {
	name       string
	Path       string
	Type       reflect.Type
	Properties map[string]int
	dir        string
	cssFile    string
	jsFile     string
}

func (k *Klass) Name() string {
	return k.Path + "/" + k.name
}

func (k *Klass) NewInstance() Component {
	//log.PrintDebug("Klass:NewInstance: %q", k.Type.PkgPath()+"/"+k.Type.Name())
	v := reflect.New(k.Type)
	vi := v.Interface()
	cc, _ := vi.(Component)
	ctx := cc.Context()
	ctx.klass = k
	ctx.self = cc
	return cc
}

func makeKlass(c Component, dir string) (*Klass, error) {
	// need lock?
	componentsLock.Lock()
	defer componentsLock.Unlock()
	pct := reflect.TypeOf(c)
	if pct.Kind() != reflect.Ptr {
		return nil, errors.New("RegisterComponent: requiered pointer")
	}
	ct := pct.Elem()
	if ct.Kind() != reflect.Struct {
		return nil, errors.New("RegisterComponent: requiered pointer of struct")
	}
	p := ct.PkgPath() + "/" + ct.Name()
	if k, ok := componentKlasses[p]; ok {
		return k, fmt.Errorf("RegisterComponent: already registerd Component: %q", p)
	}
	k := &Klass{
		name: ct.Name(),
		Path: ct.PkgPath(),
		Type: ct,
		dir:  dir,
	}
	fn := ct.NumField()
	for i := 0; i < fn; i++ {
		f := ct.Field(i)
		if ft, ok := f.Tag.Lookup("exciton"); ok {
			if k.Properties == nil {
				k.Properties = make(map[string]int)
			}
			k.Properties[ft] = i
		}
	}
	componentKlasses[p] = k
	return k, nil
}

func deleteKlass(k *Klass) {
	componentsLock.Lock()
	defer componentsLock.Unlock()
	n := k.Name()
	if _, ok := componentKlasses[n]; ok {
		delete(componentKlasses, n)
	}
}
