package markup

func Keyer(key interface{}, item RenderResult) RenderResult {
	if item == nil {
		item, _ = tag("noscript", nil)
	}
	if rr, ok := item.(keyedRenderResult); ok {
		rr.setKey(key)
	}
	return item
}
