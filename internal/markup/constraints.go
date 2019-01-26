package markup

type MarkupOrChild interface {
	isMarkupOrChild()
}

type Markup interface {
	isMarkup()
	isMarkupOrChild()
	applyToNode(b Builder, n Node, on Node)
}

type ComponentMarkup interface {
	applyToComponent(c Component)
}
