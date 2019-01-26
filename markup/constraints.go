package markup

import (
	"github.com/yossoy/exciton/internal/markup"
)

type MarkupOrChild = markup.MarkupOrChild

type Markup = markup.Markup

type ComponentMarkup = markup.ComponentMarkup

// func markupOrChildsToInternal(mm []MarkupOrChild) []markup.MarkupOrChild {
// 	if mm == nil {
// 		return nil
// 	}
// 	mmm := make([]markup.MarkupOrChild, len(mm))
// 	for i, m := range mm {
// 		mmm[i] = m
// 	}
// 	return mmm
// }
