package markup

func Keyer(key interface{}, item *RenderResult) *RenderResult {
	if item == nil {
		item = Tag("noscript")
	}
	item.key = key
	return item
}
