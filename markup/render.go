package markup

type List []MarkupOrChild

func (l List) isMarkupOrChild() {}
func (l List) isRenderResult()  {}

type renderResult struct {
	name     string
	data     string // text or namespace
	klass    *Klass
	key      interface{}
	markups  []Markup
	children []*renderResult
}

func (rr *renderResult) isTextNode() bool { return rr != nil && rr.name == "" }

type RenderResult = *renderResult

func (rr *renderResult) isMarkupOrChild() {}
func (rr *renderResult) isRenderResult()  {}
