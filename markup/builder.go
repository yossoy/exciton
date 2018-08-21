package markup

import (
	"github.com/yossoy/exciton/event"
	ievent "github.com/yossoy/exciton/internal/event"
	"github.com/yossoy/exciton/internal/object"
)

type RequestAnimationFrameHandler func()
type UpdateDiffSetHandler func(ds *DiffSet)

type Builder interface {
	RenderBody(RenderResult)
	Rerender(...Component)

	ProcRequestAnimationFrame()

	Redirect(route string)
	OnRedirect(route string)

	UserData() interface{}
	SetUserData(interface{})
}

type builder struct {
	hostPath         string
	diffSet          *DiffSet
	renderingDiffSet *DiffSet
	nestLevel        int
	hydrating        bool
	mountAll         bool
	rootNode         *node
	mounter          []Mounter
	delayUpdater     []Component
	elements         *object.ObjectMap
	components       *object.ObjectMap
	rafHandler       RequestAnimationFrameHandler
	updateHandler    UpdateDiffSetHandler
	scheduled        bool //TODO: need atomic proc?
	rootRenderResult RenderResult
	rootRenderNode   *node
	rootComponent    Component
	userData         interface{}
	keyGenerator     func() object.ObjectKey
	route            string
}

type Buildable interface {
	Builder() Builder
	EventRoot() string
}

func NewBuilder(hostPath string) Builder {
	rn := &node{}
	b := &builder{
		hostPath:         hostPath,
		diffSet:          &DiffSet{rootNode: rn},
		renderingDiffSet: &DiffSet{rootNode: rn},
		rootNode:         rn,
		elements:         object.NewObjectMap(),
		components:       object.NewObjectMap(),
		delayUpdater:     make([]Component, 0, 16),
	}
	rn.builder = b
	return b
}

func NewAsyncBuilder(hostPath string, raf RequestAnimationFrameHandler, udh UpdateDiffSetHandler) Builder {
	rn := &node{}
	b := &builder{
		hostPath:         hostPath,
		diffSet:          &DiffSet{rootNode: rn},
		renderingDiffSet: &DiffSet{rootNode: rn},
		rootNode:         rn,
		elements:         object.NewObjectMap(),
		components:       object.NewObjectMap(),
		delayUpdater:     make([]Component, 0, 16),
		rafHandler:       raf,
		updateHandler:    udh,
	}
	rn.builder = b
	return b
}

func (b *builder) RenderBody(rr RenderResult) {
	b.diffSet.reset()

	b.rootRenderResult = rr
	b.rootRenderNode = diff(b, nil, rr, b.rootNode, false)
	if b.rootRenderNode.component != nil {
		b.rootComponent = b.rootRenderNode.component
	}
	b.rerender()
	if b.rafHandler != nil && !b.scheduled {
		b.scheduled = true
		b.rafHandler()
	}
}

func (b *builder) ProcRequestAnimationFrame() {
	if b.rafHandler == nil || b.updateHandler == nil {
		return
	}
	b.scheduled = false
	b.rerender()
	if b.diffSet.hasDiff() {
		b.diffSet, b.renderingDiffSet = b.renderingDiffSet, b.diffSet
		b.updateHandler(b.renderingDiffSet)
		b.diffSet.reset()
	}
}

func (b *builder) Rerender(c ...Component) {
	if len(c) == 0 {
		b.delayUpdater = b.delayUpdater[0:0]
		b.delayUpdater = append(b.delayUpdater, nil)
	} else {
		for _, cc := range c {
			b.enqueueRender(cc)
		}
	}
	if b.rafHandler != nil {
		if !b.scheduled {
			b.scheduled = true
			b.rafHandler()
		}
		return
	}
	b.rerender()
}

func (b *builder) enqueueRender(c Component) {
	if len(b.delayUpdater) == 1 && b.delayUpdater[0] == nil {
		return
	}
	ctx := c.Context()
	if !ctx.dirty {
		ctx.dirty = true
		b.delayUpdater = append(b.delayUpdater, c)
	}
}

func (b *builder) rerender() {
	for {
		var items []Component
		items, b.delayUpdater = b.delayUpdater, nil
		if len(items) == 0 {
			return
		}
		for _, c := range items {
			if c == nil {
				b.rootRenderNode = diff(b, b.rootRenderNode, b.rootRenderResult, b.rootNode, false)
				continue
			}
			ctx := c.Context()
			if ctx.dirty {
				renderComponent(b, c, renderOptSync, false)
				if c == b.rootComponent {
					if c.Context().base != nil {
						b.rootRenderNode = c.Context().base.(*node)
					} else {
						b.rootRenderNode = nil
					}
				}
			}
		}
	}
}

func (b *builder) setNodeValue(n *node, text string) {
	b.diffSet.setNodeValue(n, text)
	n.text = text
}

func (b *builder) addElement(n *node) {
	if n.uuid == "" {
		if b.keyGenerator != nil {
			n.uuid = b.keyGenerator()
		} else {
			n.uuid = b.elements.NewKey()
		}
		b.elements.Put(n.uuid, n)
		b.diffSet.setNodeUUID(n, n.uuid)
	}
}

func (b *builder) deleteElement(n *node) {
	if n.uuid != "" {
		b.elements.Delete(n.uuid)
		if n.component != nil {
			idmap := n.component.Context().idmap
			for k, v := range idmap {
				if v == n.uuid {
					delete(idmap, k)
					break
				}
			}
		}
		n.uuid = ""
	}
}

func (b *builder) createNode(v *tagRenderResult) *node {
	n := &node{
		tag: v.name,
	}
	if v.data != "" {
		n.ns = v.data
		b.diffSet.createNodeWithNS(n)
	} else {
		b.diffSet.createNode(n)
	}
	b.addElement(n)
	return n
}

func (b *builder) createTextNode(text string) *node {
	n := &node{text: text}
	b.diffSet.createTextNode(n)
	return n
}

func (b *builder) removeChild(p *node, c *node) {
	b.diffSet.RemoveChild(p, c)
	p.removeChild(c)
	b.deleteElement(c)
}

func (b *builder) replaceChild(p *node, e *node, d *node) /**node*/ {
	b.diffSet.ReplaceChild(p, e, d)
	r := p.replaceChild(e, d)
	//return r
	b.deleteElement(r)
}

func (b *builder) appendChild(p *node, c *node) {
	b.diffSet.appendChild(p, c)
	p.appendChild(c)
}

func (b *builder) insertBefore(t *node, c *node, pos *node) {
	b.diffSet.insertBefore(t, c, pos)
	t.insertBefore(c, pos)
}

func (b *builder) mountComponent(c Component) {
	ctx := c.Context()
	var cb *node
	if ctx.base != nil {
		cb = ctx.base.(*node)
	}
	b.diffSet.addMountComponent(cb, c)
	b.components.Put(ctx.id, c)
	if m, ok := c.(Mounter); ok {
		b.mounter = append(b.mounter, m)
	}
}

func (b *builder) unmountComponent(c Component) {
	//	if (options.beforeUnmount) options.beforeUnmount(component);
	ctx := c.Context()
	base := ctx.base
	ctx.disabled = true

	if um, ok := c.(Unmounter); ok {
		um.Unmount()
	}

	b.diffSet.addUnmountComponent(c)
	b.components.Delete(ctx.id)

	ctx.base = nil

	inner := ctx.childComponent
	if inner != nil {
		b.unmountComponent(inner)
	} else if base != nil {
		bn := base.(*node)
		//swap order?
		b.removeNode(bn)
		b.removeChildren(bn)
	}
}

func (b *builder) removeNode(n *node) {
	if p := n.parent; p != nil {
		b.removeChild(p, n)
	} else {
		b.deleteElement(n)
	}
}

func (b *builder) recollectNodeTree(n *node, unmountOnly bool) {
	if c := n.component; c != nil {
		// if node is owned by a Component, unmount that component (ends up recursing back here)
		b.unmountComponent(c)
		return
	}
	// If the node's VNode had a ref function, invoke it with null here.
	// (this is part of the React spec, and smart for unsetting references)
	//if (node[ATTR_KEY]!=null && node[ATTR_KEY].ref) node[ATTR_KEY].ref(null);
	if !unmountOnly /* || node[ATTR_KEY] == null*/ {
		b.removeNode(n)
	}
	b.removeChildren(n)
	//swap order?
}

func (b *builder) removeChildren(n *node) {
	nn := n.lastChild()
	for nn != nil {
		next := nn.previousSibling()
		b.recollectNodeTree(nn, true)
		nn = next
	}
}

func (b *builder) flushMount() {
	for _, c := range b.mounter {
		if m, ok := c.(Mounter); ok {
			m.Mount()
		}
	}
	b.mounter = b.mounter[:0]
}

func (b *builder) OnRedirect(route string) {
	// validate path
	if b.route != route {
		b.route = route
		b.Rerender()
	}
}

func (b *builder) Redirect(route string) {
	ievent.Emit(b.hostPath+"/redirectTo", event.NewValue(route))
}

func (b *builder) UserData() interface{} {
	return b.userData
}

func (b *builder) SetUserData(data interface{}) {
	b.userData = data
}
