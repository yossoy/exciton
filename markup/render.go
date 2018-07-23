package markup

type List []MarkupOrChild

func (l List) isMarkupOrChild() {}
func (l List) isRenderResult()  {}

type RenderResult interface {
	MarkupOrChild
	isTextNode() bool
	isRenderResult()
}

type keyedRenderResult interface {
	Key() interface{}
	setKey(interface{})
}

type textRenderResult struct {
	text string
}

func (rr *textRenderResult) isTextNode() bool { return true }
func (rr *textRenderResult) isMarkupOrChild() {}
func (rr *textRenderResult) isRenderResult()  {}

type tagRenderResult struct {
	name     string
	data     string // text or namespace
	key      interface{}
	markups  []Markup
	children []RenderResult
}

func (rr *tagRenderResult) isTextNode() bool     { return false }
func (rr *tagRenderResult) isMarkupOrChild()     {}
func (rr *tagRenderResult) isRenderResult()      {}
func (rr *tagRenderResult) Key() interface{}     { return rr.key }
func (rr *tagRenderResult) setKey(k interface{}) { rr.key = k }

type componentRenderResult struct {
	name     string
	klass    *Klass
	key      interface{}
	markups  []Markup
	children []RenderResult
}

func (rr *componentRenderResult) isTextNode() bool     { return false }
func (rr *componentRenderResult) isMarkupOrChild()     {}
func (rr *componentRenderResult) isRenderResult()      {}
func (rr *componentRenderResult) Key() interface{}     { return rr.key }
func (rr *componentRenderResult) setKey(k interface{}) { rr.key = k }

type delayRenderResult struct {
	data    interface{}
	proc    func(b *Builder) RenderResult
	compare func(n *node, hydrating bool) bool
}

func (rr *delayRenderResult) isTextNode() bool { return false }
func (rr *delayRenderResult) isMarkupOrChild() {}
func (rr *delayRenderResult) isRenderResult()  {}
