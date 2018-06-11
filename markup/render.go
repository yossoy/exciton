package markup

type List []MarkupOrChild

func (l List) isMarkupOrChild() {}
func (l List) isRenderResult()  {}

type RenderResult struct {
	name     string
	data     string // text or namespace
	klass    *Klass
	key      interface{}
	markups  []Markup
	children []*RenderResult
}

func (rr *RenderResult) isTextNode() bool { return rr != nil && rr.name == "" }

//type RenderResult *renderResult

func (rr *RenderResult) isMarkupOrChild() {}
func (rr *RenderResult) isRenderResult()  {}
