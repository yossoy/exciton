package markup

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type klassPathInfo struct {
	pkgPath string
	klasses map[string]*Klass
	id      string
	dir     string
}

var (
	componentsLock = sync.RWMutex{}
	klassPaths     = make(map[string]*klassPathInfo)
	klassPathIDs   = make(map[string]*klassPathInfo)
)

type Klass struct {
	name       string
	pathInfo   *klassPathInfo
	Type       reflect.Type
	Properties map[string]int
	cssFile    string
	jsFile     string
}

func (k *Klass) Name() string {
	return k.pathInfo.pkgPath + "/" + k.name
}

func (k *Klass) ClassName() string {
	return escapeClassName(k.Name())
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
	var kpi *klassPathInfo
	var ok bool
	if kpi, ok = klassPaths[ct.PkgPath()]; ok {
		if k, ok := kpi.klasses[ct.Name()]; ok {
			return k, fmt.Errorf("RegisterComponent: already registerd Component: %q", ct.PkgPath()+"/"+ct.Name())
		}
		if kpi.dir != dir {
			return nil, fmt.Errorf("RegisterComponent: invalid caller path: %q", dir)
		}
	} else {
		kpid := fmt.Sprintf("eXcItOnCoMpOnEnT_%d", len(klassPathIDs))
		kpi = &klassPathInfo{
			klasses: make(map[string]*Klass),
			id:      kpid,
			pkgPath: ct.PkgPath(),
			dir:     dir,
		}
		klassPaths[ct.PkgPath()] = kpi
		klassPathIDs[kpid] = kpi
	}
	k := &Klass{
		name:     ct.Name(),
		pathInfo: kpi,
		Type:     ct,
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
	kpi.klasses[ct.Name()] = k
	return k, nil
}

func deleteKlass(k *Klass) {
	componentsLock.Lock()
	defer componentsLock.Unlock()
	if kb, ok := klassPaths[k.pathInfo.pkgPath]; ok {
		if _, ok = kb.klasses[k.name]; ok {
			delete(kb.klasses, k.name)
		}
		if len(kb.klasses) == 0 {
			delete(klassPaths, k.pathInfo.pkgPath)
			delete(klassPathIDs, k.pathInfo.id)
		}
	}
}
