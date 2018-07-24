package svg

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/yossoy/exciton/log"

	"github.com/yossoy/exciton/markup"
)

//go:generate go run gen_elem.go
//go:generate go run gen_event.go

type elemType int

const (
	elemTypeElem elemType = iota
	elemTypeChar
)

type svgElementStackItem struct {
	etype    elemType
	elem     xml.Token
	children []*svgElementStackItem
}

func renderSub(e *svgElementStackItem, parentNS string) markup.RenderResult {
	if e.etype == elemTypeChar {
		return markup.Text(string(e.elem.(xml.CharData)))
	}
	t := e.elem.(xml.StartElement)
	markups := make([]markup.MarkupOrChild, 0, len(e.children)+len(t.Attr))
	tns := t.Name.Space
	for _, a := range t.Attr {
		ans := a.Name.Space
		if ans == "xmlns" || a.Name.Local == "xmlns" {
			continue
		}
		var am markup.MarkupOrChild
		if ans == "" || ans == tns {
			am = markup.Attribute(a.Name.Local, a.Value)
		} else {
			am = markup.AttributeNS(a.Name.Local, ans, a.Value)
		}
		markups = append(markups, am)
	}
	for _, c := range e.children {
		markups = append(markups, renderSub(c, tns))
	}
	return markup.TagWithNS(t.Name.Local, tns, markups...)
}

func SVGtoRenderResult(r io.Reader) (markup.RenderResult, error) {
	//TODO: proxy support
	dec := xml.NewDecoder(r)
	var stack []*svgElementStackItem
	var top *svgElementStackItem

	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch vt := t.(type) {
		case xml.StartElement:
			s := &svgElementStackItem{
				etype: elemTypeElem,
				elem:  vt.Copy(),
			}
			stack = append(stack, s)
		case xml.EndElement:
			t := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				top = t
			} else {
				l := stack[len(stack)-1]
				l.children = append(l.children, t)
			}
		case xml.CharData:
			if len(stack) == 0 {
				// skip?
				break
			}
			t := stack[len(stack)-1]
			s := &svgElementStackItem{
				etype: elemTypeChar,
				elem:  vt.Copy(),
			}
			t.children = append(t.children, s)
		case xml.ProcInst:
		case xml.Comment:
		case xml.Directive:
			log.PrintDebug("Directive!: %q", string(vt))
		default:
			return nil, fmt.Errorf("invalid token type: %v", t)
		}
	}
	if top == nil {
		return nil, fmt.Errorf("parse error")
	}
	return renderSub(top, ""), nil
}
