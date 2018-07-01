package markup

import (
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
)

type Core struct {
	klass           *Klass
	self            Component
	childComponent  Component
	parentComponent Component
	children        []*RenderResult
	base            *node
	disabled        bool
	builder         *Builder
	dirty           bool
	key             interface{}
	parentMarkups   []Markup
}

type renderOptType int

const (
	renderOptNone renderOptType = iota
	renderOptSync
	renderOptAsync
)

type Component interface {
	Context() *Core
	Builder() *Builder
	Key() interface{}
	Render() *RenderResult
	Classes(...string) MarkupOrChild
}

func (c *Core) Context() *Core            { return c }
func (c *Core) Key() interface{}          { return c.key }
func (c *Core) Builder() *Builder         { return c.builder }
func (c *Core) Children() []*RenderResult { return c.children }

func (c *Core) Classes(classes ...string) MarkupOrChild {
	k := c.klass
	if k.cssFile == "" {
		return Classes(classes...)
	}
	ccs := make(classApplyer, len(classes))
	prefix := k.pathInfo.id + "-" + strings.TrimSuffix(k.cssFile, filepath.Ext(k.cssFile)) + "-"
	for i, class := range classes {
		ccs[i] = prefix + class
	}
	return ccs
}

type ComponentInstance func(m ...MarkupOrChild) *RenderResult

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

type WillMounter interface {
	WillMount()
}

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

func WithGlobalCSS() ComponentRegisterParameter {
	return func(k *Klass) error {
		k.cssIsGlobal = true
		return nil
	}
}

func WithComponentStyleSheet(css string) ComponentRegisterParameter {
	return func(k *Klass) error {
		if k.cssFile != "" {
			return fmt.Errorf("component %q already has component css file (%q)", k.Name(), k.cssFile)
		}
		k.cssFile = css
		return nil
	}
}

func WithComponentScript(js string) ComponentRegisterParameter {
	return func(k *Klass) error {
		if k.jsFile != "" {
			return fmt.Errorf("component %q already has component js file (%q)", k.Name(), k.jsFile)
		}
		k.jsFile = js
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

func registerComponent(c Component, dir string, params []ComponentRegisterParameter) (ComponentInstance, *Klass, error) {
	k, err := makeKlass(c, dir)
	if err != nil {
		return nil, nil, err
	}
	for _, p := range params {
		if err := p(k); err != nil {
			return nil, nil, err
		}
	}
	return ComponentInstance(func(m ...MarkupOrChild) *RenderResult {
		markups, children, err := splitMarkupOrChild(m)
		if err != nil {
			panic(err)
		}
		children2, err := flattenChildren(children)
		if err != nil {
			panic(err)
		}
		rr := &RenderResult{
			name:     k.Name(),
			markups:  markups,
			children: children2,
			klass:    k,
		}
		return rr
	}), k, nil
}

func RegisterComponent(c Component, params ...ComponentRegisterParameter) (ComponentInstance, *Klass, error) {
	_, fp, _, ok := runtime.Caller(1)
	if !ok {
		return nil, nil, fmt.Errorf("invalid caller")
	}
	return registerComponent(c, filepath.Dir(fp), params)
}

func MustRegisterComponent(c Component, params ...ComponentRegisterParameter) (ComponentInstance, *Klass) {
	_, fp, _, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Errorf("invalid caller"))
	}
	ci, k, err := registerComponent(c, filepath.Dir(fp), params)
	if err != nil {
		panic(err)
	}
	return ci, k
}

func unregisterComponent(ci ComponentInstance) {
	rr := ci()
	deleteKlass(rr.klass)
}

func createComponent(b *Builder, vnode *RenderResult) Component {
	c := vnode.klass.NewInstance()
	c.Context().builder = b

	// call initializer
	if i, ok := c.(Initializer); ok {
		i.Initialize()
	}
	return c
}

func renderComponent(b *Builder, c Component, renderOpt renderOptType, isChild bool) {
	ctx := c.Context()
	initialBase := ctx.base
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

		if rendered != nil {
			if ctx.parentMarkups != nil {
				rendered.markups = append(rendered.markups, ctx.parentMarkups...)
			}
		}

		if rendered != nil && rendered.klass != nil {
			inst = initialChildComponent

			if inst != nil && inst.Context().klass == rendered.klass && inst.Key() == rendered.key {
				setComponentProps(b, inst, renderOptSync, rendered.markups, rendered.children)
			} else {
				toUnmount = inst
				inst = createComponent(b, rendered)
				ctx.childComponent = inst
				inst.Context().parentComponent = c
				setComponentProps(b, inst, renderOptNone, rendered.markups, rendered.children)
				renderComponent(b, inst, renderOptSync, true)
			}
			base = inst.Context().base
		} else {
			cbase := ctx.base
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
		if m, ok := c.(Mounter); ok {
			b.mounter = append(b.mounter, m)
		}
	} else if !skip {
		if du, ok := c.(DidUpdate); ok {
			du.DidUpdate()
		}
	}
	if b.nestLevel == 0 && !isChild {
		b.flushMount()
	}
}

func setComponentProps(b *Builder, c Component, renderOpt renderOptType, markups []Markup, children []*RenderResult) {
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

func buildComponentFromVNode(b *Builder, dom *node, vnode *RenderResult) *node {
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
		dom = c.Context().base
	} else {
		if origComponent != nil && !isDirectOwner {
			b.unmountComponent(origComponent)
			oldDom = nil
			dom = nil
		}

		c = createComponent(b, vnode)
		if dom != nil {
			oldDom = nil
		}
		setComponentProps(b, c, renderOptSync, vnode.markups, vnode.children)
		dom = c.Context().base

		if oldDom != nil && dom != oldDom {
			oldDom.component = nil
			b.recollectNodeTree(oldDom, false)
		}
	}

	return dom
}
