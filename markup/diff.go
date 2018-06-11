package markup

func diff(b *Builder, dom *node, vnode *RenderResult, parent *node, componentRoot bool) *node {
	if b.nestLevel == 0 {
		if dom != nil && dom.uuid == "" {
			b.hydrating = true
		} else {
			b.hydrating = false
		}

	}
	b.nestLevel++

	ret := idiff(b, dom, vnode, componentRoot)

	if parent != nil && ret.parent != parent {
		b.appendChild(parent, ret)
	}

	b.nestLevel--
	if b.nestLevel == 0 {
		b.hydrating = false
		// invoke queued componentDidMount lifecycle methods
		if !componentRoot {
			b.flushMount()
		}
	}

	return ret
}

func idiff(b *Builder, dom *node, vnode *RenderResult, componentRoot bool) *node {
	out := dom
	//prevSvgMode = isSvgMode;

	// empty values (null, undefined, booleans) render as empty Text nodes
	//if (vnode==null || typeof vnode==='boolean') vnode = '';
	if vnode == nil {
		vnode = Tag("noscript")
	}

	// Fast case: Strings & Numbers create/update Text nodes.
	if vnode.isTextNode() {
		if dom.isTextNode() && dom.parent != nil && (dom.component == nil || componentRoot) {
			if dom.text != vnode.data {
				b.setNodeValue(dom, vnode.data)
			}
		} else {
			// it wasn't a Text node: replace it with one and recycle the old Element
			out = b.createTextNode(vnode.data)
			if dom != nil {
				if dom.parent != nil {
					b.replaceChild(dom.parent, out, dom)
				}
				b.recollectNodeTree(dom, true)
			}
		}

		//out[ATTR_KEY] = true;

		return out
	}

	// If the VNode represents a Component, perform a component diff:
	if vnode.klass != nil {
		return buildComponentFromVNode(b, dom, vnode)
	}

	// Tracks entering and exiting SVG namespace when descending through the tree.
	//isSvgMode = vnodeName==='svg' ? true : vnodeName==='foreignObject' ? false : isSvgMode;

	if dom == nil || dom.tag != vnode.name {
		out = b.createNode(vnode)

		if dom != nil {
			// move children into the replacement node
			for dom.firstChild() != nil {
				b.appendChild(out, dom.firstChild())
			}

			// if the previous Element was mounted into the DOM, replace it inline
			if dom.parent != nil {
				b.replaceChild(dom.parent, out, dom)
			}

			// recycle the old element (skips non-Element node types)
			b.recollectNodeTree(dom, true)
		}
	}

	vchildren := vnode.children
	fc := out.firstChild()
	if !b.hydrating && len(vchildren) == 1 && vchildren[0].isTextNode() && fc != nil && fc.isTextNode() && fc.nextSibling() == nil {
		// Optimization: fast-path for elements containing a single TextNode:
		if fc.text != vchildren[0].data {
			b.setNodeValue(fc, vchildren[0].data)
		}
	} else if len(vchildren) > 0 || fc != nil {
		// otherwise, if there are existing or new children, diff them:
		innerDiffNode(
			b,
			out,
			vchildren,
			b.hydrating /* || props.dangerouslySetInnerHTML != null*/)
	}

	// Apply attributes/props from VNode to the DOM Element:
	diffMarkups(b, out, vnode.markups)

	// restore previous SVG mode: (in case we're exiting an SVG namespace)
	//TODO:
	//isSvgMode = prevSvgMode;

	return out
}

func innerDiffNode(b *Builder, dom *node, vchildren []*RenderResult, isHydrating bool) {
	//originalChildren := dom.children
	children := make([]*node, 0, len(dom.children))
	keyed := make(map[interface{}]*node)
	min := 0

	// Build up a map of keyed children and an Array of unkeyed children:
	for _, child := range dom.children {
		var key interface{}
		if len(vchildren) > 0 {
			key = child.key
		}
		if key != nil {
			keyed[key] = child
		} else {
			children = append(children, child)
		}
	}

	for i, vchild := range vchildren {
		var child *node
		key := vchild.key
		// attempt to find a node based on key matching
		if key != nil {
			if c, ok := keyed[key]; ok {
				child = c
				delete(keyed, key)
			}
		} else if child == nil && min < len(children) {
			// attempt to pluck a node of the same type from the existing children
			for j := min; j < len(children); j++ {
				c := children[j]
				if c != nil && isSameNodeType(c, vchild, isHydrating) {
					child = c
					children[j] = nil
					if j == (len(children) - 1) {
						children = children[:j]
					}
					if j == min {
						min++
						// same performance?
						//children = children[1:]
					}
					break
				}
			}
		}

		// morph the matched/found/created DOM child to match vchild (deep)
		child = idiff(b, child, vchild, false)
		child.key = key

		var f *node
		if i < len(dom.children) {
			f = dom.children[i]
		}
		if child != nil && child != dom && child != f {
			if f == nil {
				b.appendChild(dom, child)
			} else if child == f.nextSibling() {
				b.removeNode(f)
			} else {
				b.insertBefore(dom, child, f)
			}
		}
	}

	// remove unused keyed children:
	for _, v := range keyed {
		if v != nil {
			b.recollectNodeTree(v, false)
		}
	}
	// remove orphaned unkeyed children:
	for i := len(children) - 1; i >= min; i-- {
		if c := children[i]; c != nil {
			b.recollectNodeTree(c, false)
		}
	}
}

func isSameNodeType(n *node, vnode *RenderResult, hydrating bool) bool {
	if vnode.isTextNode() {
		return n.isTextNode()
	}
	if vnode.klass == nil {
		return n.component == nil &&
			n.tag == vnode.name &&
			n.ns == vnode.data
	}
	if hydrating {
		return true
	}
	if n.component != nil && n.component.Context().klass == vnode.klass {
		return true
	}
	return false
}

func diffMarkups(b *Builder, dom *node, markups []Markup) {
	//TODO: メモリ確保を減らす
	//TODO: こんな複雑な事やる必要ある?
	on := node{}
	dom.properties, on.properties = on.properties, dom.properties
	dom.attributes, on.attributes = on.attributes, dom.attributes
	dom.eventListeners, on.eventListeners = on.eventListeners, dom.eventListeners
	dom.dataset, on.dataset = on.dataset, dom.dataset
	dom.classes, on.classes = on.classes, dom.classes
	dom.styles, on.styles = on.styles, dom.styles
	dom.innerHTML, on.innerHTML = on.innerHTML, dom.innerHTML
	for _, m := range markups {
		m.applyToNode(b, dom, &on)
	}
	for k := range on.attributes {
		b.diffSet.DelAttribute(dom, k)
	}
	for k := range on.properties {
		b.diffSet.delProperty(dom, k)
	}
	for k := range on.eventListeners {
		b.diffSet.RemoveEventListener(dom, k)
	}
	for k := range on.dataset {
		b.diffSet.DelDataSet(dom, k)
	}
	for k := range on.classes {
		b.diffSet.DelClassList(dom, k)
	}
	for k := range on.styles {
		b.diffSet.DelStyle(dom, k)
	}
	if on.innerHTML != "" && dom.innerHTML == "" {
		b.diffSet.AddInnerHTML(dom, "")
	}
}
