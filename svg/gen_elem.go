// +build ignore

package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// elemNameMap translates lowercase HTML tag names from the MDN source into a
// proper Go style name with MixedCaps and initialisms:
//
//  https://github.com/golang/go/wiki/CodeReviewComments#mixed-caps
//  https://github.com/golang/go/wiki/CodeReviewComments#initialisms
//
var elemNameMap = map[string]string{
	"a": "Anchor",
}

func main() {
	doc, err := goquery.NewDocument("https://developer.mozilla.org/en-US/docs/Web/SVG/Element")
	if err != nil {
		panic(err)
	}

	file, err := os.Create("elem.gen.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprint(file, `// Package svg defines markup to create SVG elements.
//
// Generated from "SVG element reference" by Mozilla Contributors,
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element, licensed under
// CC-BY-SA 2.5.
package svg

import mkup "github.com/yossoy/exciton/markup"
`)

	doc.Find(".quick-links a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		if !strings.HasPrefix(link, "/en-US/docs/Web/SVG/Element/") {
			return
		}

		if s.Parent().Find(".icon-trash, .icon-thumbs-down-alt, .icon-warning-sign").Length() > 0 {
			return
		}

		desc, _ := s.Attr("title")

		text := s.Text()

		name := text[1 : len(text)-1]

		writeElem(file, name, desc, link)
	})
}

func writeElem(w io.Writer, name, desc, link string) {
	funName := elemNameMap[name]
	if funName == "" {
		funName = capitalize(name)
	}

	// Descriptions for elements generally read as:
	//
	//  The SVG <foobar> element ...
	//
	// Because these are consistent (sometimes with varying captalization,
	// however) we can exploit that fact to reword the documentation in proper
	// Go style:
	//
	//  Foobar ...
	//
	generalLowercase := fmt.Sprintf("the <%s> svg element", strings.ToLower(name))

	// Replace a number of 'no-break space' unicode characters which exist in
	// the descriptions with normal spaces.
	desc = strings.Replace(desc, "\u00a0", " ", -1)
	if l := len(generalLowercase); len(desc) > l && strings.HasPrefix(strings.ToLower(desc), generalLowercase) {
		desc = fmt.Sprintf("%s%s", funName, desc[l:])
	}

	fmt.Fprintf(w, `%s
//
// https://developer.mozilla.org%s
func %s(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.TagWithNS("%s", mkup.SVGNamespace, markup...)
}
`, descToComments(desc), link, funName, name)
}

func capitalize(s string) string {
	ss := strings.Split(s, "-")
	for i, sss := range ss {
		ss[i] = strings.ToUpper(sss[:1]) + sss[1:]
	}
	return strings.Join(ss, "")
}

func descToComments(desc string) string {
	c := ""
	length := 80
	for _, word := range strings.Fields(desc) {
		if length+len(word)+1 > 80 {
			length = 3
			c += "\n//"
		}
		c += " " + word
		length += len(word) + 1
	}
	return c
}
