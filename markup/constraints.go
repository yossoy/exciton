package markup

type MarkupOrChild interface {
	isMarkupOrChild()
}

type ComponentMarkupOrChild interface {
	isComponentMarkupOrChild()
	isMarkupOrChild()
}

type Markup interface {
	isMarkup()
	isMarkupOrChild()
	applyToNode(b *Builder, n *node, on *node)
}

type ComponentMarkup interface {
	isComponentMarkupOrChild()
	isComponentMarkup()
	applyToComponent(c Component)
}
