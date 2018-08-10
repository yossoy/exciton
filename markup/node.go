package markup

import (
	"github.com/yossoy/exciton/internal/object"
)

type Node interface {
	NodeName() string
}

type Element interface {
	Node
}

type elemBody struct {
	n *node
}

type TextNode interface {
	Node
}

type textNodeBody struct {
	n *node
}

type node struct {
	tag, ns, text, innerHTML string
	parent                   *node
	component                Component
	key                      interface{}
	children                 []*node
	// markups
	classes                map[string]struct{}
	styles, dataset        map[string]string
	properties, attributes map[string]interface{}
	eventListeners         map[string]*EventListener
	index                  int
	uuid                   object.ObjectKey
	rootNode               bool
}

func (n *node) NodeName() string {
	return n.tag
}

func (n *node) isMount() bool {
	if n == nil {
		return false
	}
	if n.rootNode {
		return true
	}
	return n.parent.isMount()
}

func (n *node) indexPath(rootNode *node) []int {
	var r []int
	if n == nil || n == rootNode {
		return r
	}
	if n.parent != nil {
		r = n.parent.indexPath(rootNode)
	}
	return append(r, n.index)
}

func (n *node) isTextNode() bool {
	return n != nil && n.tag == ""
}

func (n *node) firstChild() *node {
	if len(n.children) == 0 {
		return nil
	}
	return n.children[0]
}

func (n *node) lastChild() *node {
	if len(n.children) == 0 {
		return nil
	}
	return n.children[len(n.children)-1]
}

func (n *node) nextSibling() *node {
	if n.parent == nil {
		panic("exciton: unmounted node")
	}
	if (n.index + 1) < len(n.parent.children) {
		return n.parent.children[n.index+1]
	}
	return nil
}

func (n *node) previousSibling() *node {
	if n.parent == nil {
		panic("exciton: unmounted node")
	}
	if 0 < n.index {
		return n.parent.children[n.index-1]
	}
	return nil
}

func (n *node) appendChild(c *node) {
	if c.parent != nil {
		c.parent.removeChild(c)
	}
	if n.children == nil {
		n.children = make([]*node, 0, 16)
	}
	c.index = len(n.children)
	c.parent = n
	n.children = append(n.children, c)
}

func (n *node) insertBefore(c *node, pos *node) {
	if pos.parent != n {
		panic("invalid pos")
	}
	c.parent = n
	idx := pos.index
	c.index = idx
	n.children = append(n.children, c)
	copy(n.children[idx+1:], n.children[idx:])
	n.children[idx] = c
	//idx = idx + 1
	for idx = idx + 1; idx < len(n.children); idx++ {
		n.children[idx].index = idx
	}
}

func (n *node) replaceChild(e, d *node) *node {
	if d.parent != n {
		panic("exciton: invalid child")
	}
	if e.parent != nil {
		e.parent.removeChild(e)
	}
	n.children[d.index] = e
	e.index = d.index
	e.parent = n
	d.parent = nil
	d.index = -1
	return d
}

func (n *node) removeChild(c *node) *node {
	if c.parent != n {
		panic("exciton: invalid child")
	}
	n.children = append(n.children[:c.index], n.children[c.index+1:]...)
	for i := c.index; i < len(n.children); i++ {
		n.children[i].index = i
	}
	c.parent = nil
	c.index = -1
	return c
}
