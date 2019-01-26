package markup

type List []MarkupOrChild

func (l List) isMarkupOrChild() {}
func (l List) isRenderResult()  {}

type RenderResult interface {
	MarkupOrChild
	isTextNode() bool
	isRenderResult()
}

type KeyedRenderResult interface {
	Key() interface{}
	SetKey(interface{})
}

type textRenderResult struct {
	text string
}

func (rr *textRenderResult) isTextNode() bool { return true }
func (rr *textRenderResult) isMarkupOrChild() {}
func (rr *textRenderResult) isRenderResult()  {}

type tagRenderResult struct {
	Name     string
	Data     string // text or namespace
	KeyValue interface{}
	Markups  []Markup
	Children []RenderResult
}

func (rr *tagRenderResult) isTextNode() bool     { return false }
func (rr *tagRenderResult) isMarkupOrChild()     {}
func (rr *tagRenderResult) isRenderResult()      {}
func (rr *tagRenderResult) Key() interface{}     { return rr.KeyValue }
func (rr *tagRenderResult) SetKey(k interface{}) { rr.KeyValue = k }

type ComponentRenderResult struct {
	Name     string
	Klass    *klass
	KeyValue interface{}
	Markups  []Markup
	Children []RenderResult
}

func (rr *ComponentRenderResult) isTextNode() bool     { return false }
func (rr *ComponentRenderResult) isMarkupOrChild()     {}
func (rr *ComponentRenderResult) isRenderResult()      {}
func (rr *ComponentRenderResult) Key() interface{}     { return rr.KeyValue }
func (rr *ComponentRenderResult) SetKey(k interface{}) { rr.KeyValue = k }

type delayRenderResult struct {
	data interface{}
	proc func(b Builder) RenderResult
}

func (rr *delayRenderResult) isTextNode() bool { return false }
func (rr *delayRenderResult) isMarkupOrChild() {}
func (rr *delayRenderResult) isRenderResult()  {}

func FuncToRenderResult(proc func(b Builder) RenderResult) *delayRenderResult {
	return &delayRenderResult{
		data: nil,
		proc: proc,
	}
}
