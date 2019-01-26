package markup

import (
	"github.com/yossoy/exciton/internal/markup"
)

func Keyer(key interface{}, item RenderResult) RenderResult {
	if item == nil {
		item = Tag("noscript")
	}
	if rr, ok := item.(markup.KeyedRenderResult); ok {
		rr.SetKey(key)
	}
	return item
}
