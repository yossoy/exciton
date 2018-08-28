package markup

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
	ievent "github.com/yossoy/exciton/internal/event"
	"github.com/yossoy/exciton/internal/object"
	"github.com/yossoy/exciton/log"
)

type Core struct {
	klass           *Klass
	id              string
	self            Component
	childComponent  Component
	parentComponent Component
	children        []RenderResult
	base            Node
	disabled        bool
	builder         Builder
	dirty           bool
	key             interface{}
	parentMarkups   []Markup
	idmap           map[string]string
}

type renderOptType int

const (
	renderOptNone renderOptType = iota
	renderOptSync
	renderOptAsync
)

type Component interface {
	json.Marshaler
	Context() *Core
	Builder() Builder
	Key() interface{}
	Render() RenderResult
	Classes(...string) MarkupOrChild
	ID(string) MarkupOrChild
	GetProperty(string) (interface{}, bool)
}

func (c *Core) Context() *Core           { return c }
func (c *Core) Key() interface{}         { return c.key }
func (c *Core) Builder() Builder         { return c.builder }
func (c *Core) Children() []RenderResult { return c.children }
func (c *Core) ResourcesFilePath(fn string) string {
	return path.Join(c.klass.ResourcePathBase(), "resources", fn)
}

type componentIDApplyer struct {
	c     *Core
	idstr string
}

func (cia componentIDApplyer) isMarkup()        {}
func (cia componentIDApplyer) isMarkupOrChild() {}
func (cia componentIDApplyer) applyToNode(b Builder, n Node, on Node) {
	nn := n.(*node)
	cia.c.idmap[cia.idstr] = nn.uuid
}

func (c *Core) ID(id string) MarkupOrChild {
	return componentIDApplyer{
		c:     c,
		idstr: id,
	}
}
func (c *Core) GetProperty(name string) (interface{}, bool) {
	if c.self == nil {
		return nil, false
	}
	if idx, ok := c.klass.Properties[name]; ok {
		v := reflect.ValueOf(c.self)
		return v.Elem().Field(idx).Interface(), true
	}
	return nil, false
}

func (c *Core) Classes(classes ...string) MarkupOrChild {
	k := c.klass
	if k.localCSSFile == "" {
		return Classes(classes...)
	}
	ccs := make(classApplyer, len(classes))
	prefix := k.pathInfo.id + "-" + strings.TrimSuffix(k.localCSSFile, filepath.Ext(k.localCSSFile)) + "-"
	for i, class := range classes {
		ccs[i] = prefix + class
	}
	return ccs
}

func (c *Core) ClientJSEvent(name string, funcName string, arguments ...interface{}) MarkupOrChild {
	k := c.klass
	prefix := ""
	if k.localJSFile != "" {
		prefix = k.pathInfo.id + "-" + strings.TrimSuffix(k.localJSFile, filepath.Ext(k.localJSFile))
	}
	return &EventListener{
		Name:               name,
		clientScriptPrefix: prefix,
		scriptHandlerName:  funcName,
		scriptArguments:    arguments,
	}
}

func (c *Core) CallClientFunction(funcName string, arguments ...interface{}) (json.RawMessage, error) {
	arg := struct {
		FuncName  string        `json:"funcName"`
		Arguments []interface{} `json:"arguments"`
	}{
		FuncName:  funcName,
		Arguments: arguments,
	}
	bc := &browserCommand{
		Command:  "callClientFunction",
		Target:   c,
		Argument: &arg,
	}
	result := ievent.EmitWithResult(c.Builder().(*builder).hostPath+"/browserSync", event.NewValue(bc))
	log.PrintDebug("call result: %v", result)
	if err := result.Error(); err != nil {
		return nil, err
	}
	return result.Value().Encode()
}

func (c *Core) MarshalJSON() ([]byte, error) {
	s := struct {
		ClassID    string `json:"classId"`
		ID         string `json:"id"`
		LocalJSKey string `json:"localJSKey"`
		URLBase    string `json:"urlBase"`
	}{
		ClassID:    c.klass.Name(),
		ID:         c.id,
		LocalJSKey: c.klass.localJSKey(),
		URLBase:    c.klass.ResourcePathBase(),
	}
	return json.Marshal(&s)
}

type ComponentInstance func(m ...MarkupOrChild) RenderResult

// Mounter is an optional interface that a Component can implement in order
// to receive component mount events.
type Mounter interface {
	Mount()
}

// Unmounter is an optional interface that a Component can implement in order
// to receive component unmount events.
type Unmounter interface {
	Unmount()
}

// WillMounter is an optional interface that a Component can implement in order
// to receive component willnmount events.
type WillMounter interface {
	WillMount()
}

// ShouldUpdate is an optional interface that can be implemented to determine
// whether a component needs to be updated.
type ShouldUpdate interface {
	ShouldUpdate() bool
}

type WillUpdate interface {
	WillUpdate()
}

type DidUpdate interface {
	DidUpdate()
}

type Initializer interface {
	Initialize()
}

type ComponentRegisterParameter func(k *Klass) error

func WithGlobalStyleSheet(css string) ComponentRegisterParameter {
	return func(k *Klass) error {
		g, ok := k.pathInfo.cssFiles[css]
		if ok && !g {
			return fmt.Errorf("css %q registerd as component css file by other component", css)
		}
		if k.pathInfo.cssFiles == nil {
			k.pathInfo.cssFiles = make(map[string]bool)
		}
		k.pathInfo.cssFiles[css] = true
		return nil
	}
}

func WithComponentStyleSheet(css string) ComponentRegisterParameter {
	return func(k *Klass) error {
		g, ok := k.pathInfo.cssFiles[css]
		if ok && g {
			return fmt.Errorf("css %q registerd as global css file by other component", css)
		}
		if k.localCSSFile != "" {
			return fmt.Errorf("component %q has another component css file (%s)", k.Name(), k.localCSSFile)
		}
		if k.pathInfo.cssFiles == nil {
			k.pathInfo.cssFiles = make(map[string]bool)
		}
		k.localCSSFile = css
		k.pathInfo.cssFiles[css] = false
		return nil
	}
}

func WithGlobalScript(js string) ComponentRegisterParameter {
	return func(k *Klass) error {
		g, ok := k.pathInfo.jsFiles[js]
		if ok && !g {
			return fmt.Errorf("js %q registerd as component js file by other component", js)
		}
		if k.pathInfo.jsFiles == nil {
			k.pathInfo.jsFiles = make(map[string]bool)
		}
		k.pathInfo.jsFiles[js] = true
		return nil
	}
}

func WithComponentScript(js string) ComponentRegisterParameter {
	return func(k *Klass) error {
		g, ok := k.pathInfo.jsFiles[js]
		if ok && g {
			return fmt.Errorf("js %q already registerd as global js file by other component", js)
		}
		if k.localJSFile != "" {
			return fmt.Errorf("component %q has another component js file (%s)", k.Name(), k.localJSFile)
		}
		if k.pathInfo.jsFiles == nil {
			k.pathInfo.jsFiles = make(map[string]bool)
		}
		k.localJSFile = js
		k.pathInfo.jsFiles[js] = false
		return nil
	}
}

type InitInfo struct {
	Klass  *Klass
	timing driver.InitProcTiming
	si     *driver.StartupInfo
}

type EventHandler func(c Component, e *event.Event)

var eventGroup event.Group

func addEventHandlerSub(eventRoot event.Group, timing driver.InitProcTiming) (event.Group, error) {
	if timing == driver.InitProcTimingPreStartup {
		return nil, fmt.Errorf("cannot initialize event in InitProcTimingPreStartup: %d", timing)
	}
	if eventGroup == nil {
		var err error
		eventGroup, err = eventRoot.AddGroup("/components/:windowId/:instanceId")
		if err != nil {
			return nil, err
		}
	}
	return eventGroup, nil
}

func eventToComponent(e *event.Event) (Component, error) {
	wid := e.Params["windowId"]
	id := e.Params["instanceId"]
	wi := object.Windows.Get(wid)
	if wi == nil {
		return nil, fmt.Errorf("window not found: %s", wid)
	}
	b := wi.(Buildable).Builder().(*builder)
	ci := b.components.Get(id)
	if ci == nil {
		return nil, fmt.Errorf("component not found: %s", id)
	}
	c, ok := ci.(Component)
	if !ok {
		panic("registerd object is not component!")
	}
	return c, nil
}

func (ii *InitInfo) AddHandler(name string, handler EventHandler) error {
	group, err := addEventHandlerSub(ii.si.AppEventRoot, ii.timing)
	if err != nil {
		return err
	}
	return group.AddHandler(name, func(e *event.Event) {
		log.PrintDebug("InitInfo: handler called: %q, %v", name, e)
		c, err := eventToComponent(e)
		if err != nil {
			log.PrintError("event is not handled: %v", err)
			return
		}
		handler(c, e)
	})
}

func (ii *InitInfo) Router() driver.Router {
	return ii.si.Router
}

type ClassInitProc func(k *Klass, si *InitInfo) error

func WithClassInitializer(timing driver.InitProcTiming, proc ClassInitProc) ComponentRegisterParameter {
	return func(k *Klass) error {
		driver.AddInitProc(timing, func(si *driver.StartupInfo) error {
			ii := &InitInfo{
				Klass:  k,
				timing: timing,
				si:     si,
			}
			return proc(k, ii)
		})
		return nil
	}
}

func filePathToFileURI(path string) *url.URL {
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	url := &url.URL{}
	url.Scheme = "file"
	url.Host = ""
	url.Path = path
	return url
}

func escapeClassName(name string) string {
	sb := strings.Builder{}
	for _, c := range name {
		switch {
		case '0' <= c && c <= '9':
			fallthrough
		case 'A' <= c && c <= 'Z':
			fallthrough
		case 'a' <= c && c <= 'z':
			fallthrough
		case c == '_':
			sb.WriteRune(c)
		case c < 256:
			sb.WriteByte('\\')
			sb.WriteRune(c)
		default:
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

func registerComponent(c Component, dir string, params []ComponentRegisterParameter, klassOnly bool) (ComponentInstance, *Klass, error) {
	k, err := makeKlass(c, dir)
	if err != nil {
		return nil, nil, err
	}
	for _, p := range params {
		if err := p(k); err != nil {
			return nil, nil, err
		}
	}
	if klassOnly {
		return nil, k, nil
	}
	return ComponentInstance(func(m ...MarkupOrChild) RenderResult {
		markups, children, err := splitMarkupOrChild(m)
		if err != nil {
			panic(err)
		}
		children2, err := flattenChildren(children)
		if err != nil {
			panic(err)
		}
		rr := &componentRenderResult{
			name:     k.Name(),
			markups:  markups,
			children: children2,
			klass:    k,
		}
		return rr
	}), k, nil
}

func RegisterComponent(c Component, params ...ComponentRegisterParameter) (ComponentInstance, error) {
	_, fp, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("invalid caller")
	}
	ci, _, err := registerComponent(c, filepath.Dir(fp), params, false)
	return ci, err
}

func MustRegisterComponent(c Component, params ...ComponentRegisterParameter) ComponentInstance {
	_, fp, _, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Errorf("invalid caller"))
	}
	ci, _, err := registerComponent(c, filepath.Dir(fp), params, false)
	if err != nil {
		panic(err)
	}
	return ci
}

func RegisterKlassOnly(c Component, params ...ComponentRegisterParameter) (*Klass, error) {
	_, fp, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("invalid caller")
	}
	_, k, err := registerComponent(c, filepath.Dir(fp), params, true)
	return k, err
}

func MustRegisterKlassOnly(c Component, params ...ComponentRegisterParameter) *Klass {
	_, fp, _, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Errorf("invalid caller"))
	}
	_, k, err := registerComponent(c, filepath.Dir(fp), params, true)
	if err != nil {
		panic(err)
	}
	return k
}

func unregisterComponent(ci ComponentInstance) {
	rr := ci()
	rrr, ok := rr.(*componentRenderResult)
	if !ok {
		panic(fmt.Errorf("invalid result: %v", rr))
	}
	deleteKlass(rrr.klass)
}

func createComponent(b *builder, vnode *componentRenderResult) Component {
	c := vnode.klass.NewInstance()
	c.Context().builder = b
	c.Context().id = b.components.NewKey()

	// call initializer
	if i, ok := c.(Initializer); ok {
		i.Initialize()
	}
	return c
}

func renderComponent(b *builder, c Component, renderOpt renderOptType, isChild bool) {
	ctx := c.Context()
	var initialBase *node
	if ctx.base != nil {
		initialBase = ctx.base.(*node)
	}
	bUpdate := initialBase != nil
	skip := false
	initialChildComponent := ctx.childComponent

	if bUpdate {
		if scu, ok := c.(ShouldUpdate); ok {
			if !scu.ShouldUpdate() {
				skip = true
			} else if wu, ok := c.(WillUpdate); ok {
				wu.WillUpdate()
			}
		}
	}
	ctx.dirty = false
	if !skip {
		rendered := c.Render()
		var toUnmount Component
		var inst Component
		var base *node
		procComponentResult := false

		switch vt := rendered.(type) {
		case nil:
		case *componentRenderResult:
			if ctx.parentMarkups != nil {
				vt.markups = append(vt.markups, ctx.parentMarkups...)
			}
			inst = initialChildComponent
			if inst != nil && inst.Context().klass == vt.klass && inst.Key() == vt.key {
				setComponentProps(b, inst, renderOptSync, vt.markups, vt.children)
			} else {
				toUnmount = inst
				inst = createComponent(b, vt)
				ctx.childComponent = inst
				inst.Context().parentComponent = c
				setComponentProps(b, inst, renderOptNone, vt.markups, vt.children)
				renderComponent(b, inst, renderOptSync, true)
			}
			base = inst.Context().base.(*node)
			procComponentResult = true
		case *textRenderResult:
		case *delayRenderResult:
		case *tagRenderResult:
			if ctx.parentMarkups != nil {
				vt.markups = append(vt.markups, ctx.parentMarkups...)
			}
		default:
			panic(fmt.Errorf("type not implement!: %v", vt))
		}
		if !procComponentResult {
			var cbase *node
			if ctx.base != nil {
				cbase = ctx.base.(*node)
			}
			toUnmount = initialChildComponent
			if toUnmount != nil {
				ctx.childComponent = nil
				cbase = nil
			}
			if initialBase != nil || renderOpt == renderOptSync {
				if cbase != nil {
					cbase.component = nil
				}
				var parent *node
				if initialBase != nil {
					parent = initialBase.parent
				}
				base = diff(b, cbase, rendered, parent, true)
			}
		}
		if initialBase != nil && base != initialBase && inst != initialChildComponent {
			baseParent := initialBase.parent
			if baseParent != nil && base != baseParent {
				b.replaceChild(baseParent, base, initialBase)

				if toUnmount == nil {
					initialBase.component = nil
					b.recollectNodeTree(initialBase, false)
				}
			}
		}
		if toUnmount != nil {
			b.unmountComponent(toUnmount)
		}
		ctx.base = base
		if base != nil && !isChild {
			componentRef := c
			t := c
			for {
				t = t.Context().parentComponent
				if t == nil {
					break
				}
				componentRef = t
				componentRef.Context().base = base
			}
			base.component = componentRef
		}
	}
	if !bUpdate {
		b.mountComponent(c)
	} else if !skip {
		if du, ok := c.(DidUpdate); ok {
			du.DidUpdate()
		}
	}
	if b.nestLevel == 0 && !isChild {
		b.flushMount()
	}
}

func setComponentProps(b *builder, c Component, renderOpt renderOptType, markups []Markup, children []RenderResult) {
	ctx := c.Context()
	if ctx.disabled {
		return
	}
	ctx.disabled = true
	if ctx.base == nil {
		if wi, ok := c.(WillMounter); ok {
			wi.WillMount()
		}
	}
	ctx.disabled = false
	//TODO: async render

	// apply to ComponentMarkup
	mm := make([]Markup, 0, len(markups))
	for _, m := range markups {
		if cm, ok := m.(ComponentMarkup); ok {
			cm.applyToComponent(c)
		} else {
			mm = append(mm, m)
		}
	}
	ctx.parentMarkups = mm
	ctx.children = children

	if renderOpt != renderOptNone {
		if renderOpt == renderOptSync || ctx.base == nil {
			renderComponent(b, c, renderOptSync, false)
		} else {
			// update
			b.enqueueRender(c)
		}
	}
}

func buildComponentFromVNode(b *builder, dom *node, vnode *componentRenderResult) *node {
	var c Component
	if dom != nil {
		c = dom.component
	}
	origComponent := c
	oldDom := dom
	isDirectOwner := (c != nil && c.Context().klass == vnode.klass)
	isOwner := isDirectOwner

	for c != nil && !isOwner {
		c = c.Context().parentComponent
		if c == nil {
			break
		}
		isOwner = c.Context().klass == vnode.klass
	}
	if c != nil && isOwner && (!b.mountAll || c.Context().childComponent != nil) {
		setComponentProps(b, c, renderOptAsync, vnode.markups, vnode.children)
		if c.Context().base != nil {
			dom = c.Context().base.(*node)
		} else {
			dom = nil
		}
	} else {
		if origComponent != nil && !isDirectOwner {
			b.unmountComponent(origComponent)
			oldDom = nil
			dom = nil
		}

		c = createComponent(b, vnode)
		setComponentProps(b, c, renderOptSync, vnode.markups, vnode.children)
		if c.Context().base != nil {
			dom = c.Context().base.(*node)
		} else {
			dom = nil
		}

		if oldDom != nil && dom != oldDom {
			oldDom.component = nil
			b.recollectNodeTree(oldDom, false)
		}
	}

	return dom
}
