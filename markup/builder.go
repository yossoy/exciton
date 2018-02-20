package markup

import "github.com/yossoy/exciton/internal/object"

type RequestAnimationFrameHandler func()
type UpdateDiffSetHandler func(ds *DiffSet)

type Builder struct {
	diffSet          *DiffSet
	renderingDiffSet *DiffSet
	nestLevel        int
	hydrating        bool
	mountAll         bool
	rootNode         *node
	mounter          []Mounter
	delayUpdater     []Component
	elements         *object.ObjectMap
	rafHandler       RequestAnimationFrameHandler
	updateHandler    UpdateDiffSetHandler
	scheduled        bool
	rootRenderResult RenderResult
	rootRenderNode   *node
	rootComponent    Component
	UserData         interface{}
	keyGenerator     func() object.ObjectKey
}

type Buildable interface {
	Builder() *Builder
}

func NewBuilder() *Builder {
	rn := &node{}
	b := &Builder{
		diffSet:          &DiffSet{rootNode: rn},
		renderingDiffSet: &DiffSet{rootNode: rn},
		rootNode:         rn,
		elements:         object.NewObjectMap(),
		delayUpdater:     make([]Component, 0, 16),
	}
	return b
}

func NewAsyncBuilder(raf RequestAnimationFrameHandler, udh UpdateDiffSetHandler) *Builder {
	rn := &node{}
	b := &Builder{
		diffSet:          &DiffSet{rootNode: rn},
		renderingDiffSet: &DiffSet{rootNode: rn},
		rootNode:         rn,
		elements:         object.NewObjectMap(),
		delayUpdater:     make([]Component, 0, 16),
		rafHandler:       raf,
		updateHandler:    udh,
	}
	return b
}

func (b *Builder) RenderBody(rr RenderResult) {
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

func (b *Builder) ProcRequestAnimationFrame() {
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

func (b *Builder) Rerender(c ...Component) {
	if len(c) == 0 {
		b.delayUpdater = b.delayUpdater[0:0]
		b.delayUpdater = append(b.delayUpdater, nil)
	} else {
		for _, cc := range c {
			b.enqueueRender(cc)
		}
	}
	if b.rafHandler != nil && !b.scheduled {
		b.scheduled = true
		b.rafHandler()
		return
	}
	b.rerender()
}

func (b *Builder) enqueueRender(c Component) {
	if len(b.delayUpdater) == 1 && b.delayUpdater[0] == nil {
		return
	}
	ctx := c.Context()
	if !ctx.dirty {
		ctx.dirty = true
		b.delayUpdater = append(b.delayUpdater, c)
	}
}

func (b *Builder) rerender() {
	for {
		items := b.delayUpdater
		b.delayUpdater = nil
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
					b.rootRenderNode = c.Context().base
				}
			}
		}
	}
}

func (b *Builder) setNodeValue(n *node, text string) {
	b.diffSet.setNodeValue(n, text)
	n.text = text
}

func (b *Builder) addElement(n *node) {
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

func (b *Builder) deleteElement(n *node) {
	if n.uuid != "" {
		b.elements.Delete(n.uuid)
		n.uuid = ""
	}
}

func (b *Builder) createNode(v RenderResult) *node {
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

func (b *Builder) createTextNode(text string) *node {
	n := &node{text: text}
	b.diffSet.createTextNode(n)
	return n
}

func (b *Builder) removeChild(p *node, c *node) {
	b.diffSet.RemoveChild(p, c)
	p.removeChild(c)
	b.deleteElement(c)
}

func (b *Builder) replaceChild(p *node, e *node, d *node) /**node*/ {
	b.diffSet.ReplaceChild(p, e, d)
	r := p.replaceChild(e, d)
	//return r
	b.deleteElement(r)
}

func (b *Builder) appendChild(p *node, c *node) {
	b.diffSet.appendChild(p, c)
	p.appendChild(c)
}

func (b *Builder) insertBefore(t *node, c *node, pos *node) {
	b.diffSet.insertBefore(t, c, pos)
	t.insertBefore(c, pos)
}

func (b *Builder) unmountComponent(c Component) {
	//	if (options.beforeUnmount) options.beforeUnmount(component);
	ctx := c.Context()
	base := ctx.base
	ctx.disabled = true

	if um, ok := c.(Unmounter); ok {
		um.Unmount()
	}

	ctx.base = nil

	inner := ctx.childComponent
	if inner != nil {
		b.unmountComponent(inner)
	} else if base != nil {
		//base.component = nil

		//swap order?
		b.removeNode(base)
		b.removeChildren(base)
	}

}

func (b *Builder) removeNode(n *node) {
	if p := n.parent; p != nil {
		b.removeChild(p, n)
	} else {
		b.deleteElement(n)
	}
}

func (b *Builder) recollectNodeTree(n *node, unmountOnly bool) {
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

func (b *Builder) removeChildren(n *node) {
	nn := n.lastChild()
	for nn != nil {
		next := nn.previousSibling()
		b.recollectNodeTree(nn, true)
		nn = next
	}
}

func (b *Builder) flushMount() {
	for _, c := range b.mounter {
		if m, ok := c.(Mounter); ok {
			m.Mount()
		}
	}
	b.mounter = b.mounter[:0]
}
