// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Event struct {
	Name      string
	EventType string
	Link      string
	Desc      string
	Category  string
}

func main() {
	// nameMap translates lowercase HTML attribute names from the MDN source
	// into a proper Go style name with MixedCaps and initialisms:
	//
	//  https://github.com/golang/go/wiki/CodeReviewComments#mixed-caps
	//  https://github.com/golang/go/wiki/CodeReviewComments#initialisms
	//
	nameMap := map[string]string{}

	// SupportedEvents
	// see domevents.go
	supportedEvents := map[string]bool{
		"SVGEvent":     true,
		"TimeEvent":    true,
		"SVGZoomEvent": true,
	}

	doc, err := goquery.NewDocument("https://developer.mozilla.org/en-US/docs/Web/Events")
	if err != nil {
		panic(err)
	}

	events := make(map[string]*Event)
	doc.Find(".standard-table").Find("tr").Each(func(i int, s *goquery.Selection) {
		cols := s.Find("td")
		if cols.Length() != 4 || cols.Find(".icon-thumbs-down-alt").Length() != 0 {
			return
		}
		spec := strings.TrimSpace(cols.Eq(2).Text())
		if !strings.Contains(spec, "SVG") {
			return
		}

		if cols.Length() == 0 || cols.Find(".icon-thumbs-down-alt").Length() != 0 {
			return
		}

		link := cols.Eq(0).Find("a").Eq(0)
		var e Event

		et := strings.TrimSpace(cols.Eq(1).Find("a").Eq(0).Text())
		if _, ok := supportedEvents[et]; !ok {
			log.Printf("Unsupported EventType(%q) has %q event.", et, link.Text())
			return
		}
		e.Name = link.Text()
		e.EventType = et
		e.Link, _ = link.Attr("href")
		e.Desc = strings.TrimSpace(cols.Eq(3).Text())
		e.Category = spec

		funName := nameMap[e.Name]
		if funName == "" {
			funName = capitalize(e.Name)
			funName = "On" + funName
		}

		if e.Desc != "" {
			e.Desc = fmt.Sprintf("%s is an event fired when %s", funName, lowercase(e.Desc))
		} else {
			e.Desc = "(no documentation)"
		}
		events[funName] = &e
	})

	var names []string
	for name := range events {
		names = append(names, name)
	}
	sort.Strings(names)

	file, err := os.Create("event.gen.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprint(file, `package svg
// Generated from "Event reference" by Mozilla Contributors,
// https://developer.mozilla.org/en-US/docs/Web/Events, licensed under CC-BY-SA 2.5.

import "github.com/yossoy/exciton/markup"
import mkup "github.com/yossoy/exciton/internal/markup"
import "github.com/yossoy/exciton/event"
`)

	for _, name := range names {
		e := events[name]
		fmt.Fprintf(file, `%s
//
// Category: %s
//
// https://developer.mozilla.org%s
func %s(listener func(e *%s)) markup.EventListener {
	return mkup.NewEventListener("%s", func(le *event.Event) {
		dispatchEventHelper%s(le, listener)
	})
}
`, descToComments(e.Desc), e.Category, e.Link[6:], name, e.EventType, e.Name, e.EventType)
	}

}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func lowercase(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
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
