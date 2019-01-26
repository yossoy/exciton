package markup

import (
	"github.com/yossoy/exciton/internal/markup"
)

func Tag(name string, mm ...MarkupOrChild) RenderResult {
	r, err := markup.Tag(name, mm)
	if err != nil {
		panic(err)
	}
	return r
}

func TagWithNS(name string, ns string, mm ...MarkupOrChild) RenderResult {
	rr, err := markup.Tag(name, mm)
	if err != nil {
		panic(err)
	}
	rr.Data = ns
	return rr
}

func Text(text string) RenderResult {
	return markup.Text(text)
}
