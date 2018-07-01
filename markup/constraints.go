package markup

type MarkupOrChild interface {
	isMarkupOrChild()
}

type Markup interface {
	isMarkup()
	isMarkupOrChild()
	applyToNode(b *Builder, n *node, on *node)
}

type ComponentMarkup interface {
	applyToComponent(c Component)
}
