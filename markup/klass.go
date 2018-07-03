package markup

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"
)

type klassPathInfo struct {
	pkgPath  string
	klasses  map[string]*Klass
	id       string
	dir      string
	jsFiles  map[string]bool
	cssFiles map[string]bool
}

var (
	componentsLock = sync.RWMutex{}
	klassPaths     = make(map[string]*klassPathInfo)
	klassPathIDs   = make(map[string]*klassPathInfo)
)

type Klass struct {
	name         string
	pathInfo     *klassPathInfo
	Type         reflect.Type
	Properties   map[string]int
	localCSSFile string
	// jsFile      string
	cssIsGlobal bool
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
	k, err := makeKlassCore(ct.PkgPath(), ct.Name(), dir)
	if err != nil {
		return nil, err
	}
	k.Type = ct
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
	return k, nil
}

func makeKlassCore(pkgPath, name, dir string) (*Klass, error) {
	var kpi *klassPathInfo
	var ok bool
	if kpi, ok = klassPaths[pkgPath]; ok {
		if k, ok := kpi.klasses[name]; ok {
			return k, fmt.Errorf("RegisterComponent: already registerd Component: %q", pkgPath+"/"+name)
		}
		if kpi.dir != dir {
			return nil, fmt.Errorf("RegisterComponent: invalid caller path: %q", dir)
		}
	} else {
		kpid := fmt.Sprintf("eXcItOnCoMpOnEnT_%d", len(klassPathIDs))
		kpi = &klassPathInfo{
			klasses: make(map[string]*Klass),
			id:      kpid,
			pkgPath: pkgPath,
			dir:     dir,
		}
		klassPaths[pkgPath] = kpi
		klassPathIDs[kpid] = kpi
	}
	k := &Klass{
		name:     name,
		pathInfo: kpi,
	}
	kpi.klasses[name] = k
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

func (k *Klass) GetResourceFile(fn string) (http.File, error) {
	return k.pathInfo.getResourceFile(fn)
}
