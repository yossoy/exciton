package markup

import (
	"fmt"
)

type tagStack struct {
	items []RenderResult
	idx   int
}

func splitMarkupOrChild(mm []MarkupOrChild) (marksups []Markup, children []RenderResult, err error) {
	for _, m := range mm {
		switch v := m.(type) {
		case nil:
		case Markup:
			marksups = append(marksups, v)
		case RenderResult:
			children = append(children, v)
		case List:
			mm, cc, e := splitMarkupOrChild(v)
			if e != nil {
				err = e
				return
			}
			marksups = append(marksups, mm...)
			children = append(children, cc...)
		default:
			err = fmt.Errorf("invalid input: %#v", m)
		}
	}
	return
}

func flattenChildren(children []RenderResult) ([]RenderResult, error) {
	if len(children) == 0 {
		return nil, nil
	}
	stack := make([]tagStack, 1, 16)
	stack = append(stack, tagStack{
		items: children,
	})
	children2 := make([]RenderResult, 0, len(children))
	for len(stack) > 0 {
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		mmm := item.items
		start := item.idx
		for idx := start; idx < len(mmm); idx++ {
			m := mmm[idx]
			children2 = append(children2, m)
		}
	}
	if len(children2) == 0 {
		return nil, nil
	}
	return children2, nil
}

func MustTag(name string, mm []MarkupOrChild) *tagRenderResult {
	t, err := Tag(name, mm)
	if err != nil {
		panic(err) // TODO: return error render result?
	}
	return t
}

func Tag(name string, mm []MarkupOrChild) (*tagRenderResult, error) {
	markups, children, err := splitMarkupOrChild(mm)
	if err != nil {
		return nil, err
	}
	children2, err := flattenChildren(children)
	if err != nil {
		return nil, err
	}
	return &tagRenderResult{
		Name:     name,
		Markups:  markups,
		Children: children2,
	}, err
}

func Text(text string) *textRenderResult {
	return &textRenderResult{
		text: text,
	}
}
