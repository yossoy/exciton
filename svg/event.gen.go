package svg
// Generated from "Event reference" by Mozilla Contributors,
// https://developer.mozilla.org/en-US/docs/Web/Events, licensed under CC-BY-SA 2.5.

import "github.com/yossoy/exciton/markup"
import mkup "github.com/yossoy/exciton/internal/markup"
import "github.com/yossoy/exciton/event"

// OnBeginEvent is an event fired when a SMIL animation element begins.
//
// Category: SVG
//
// https://developer.mozilla.org/docs/Web/Events/beginEvent
func OnBeginEvent(listener func(e *TimeEvent)) markup.EventListener {
	return mkup.NewEventListener("beginEvent", func(le *event.Event) {
		dispatchEventHelperTimeEvent(le, listener)
	})
}

// OnEndEvent is an event fired when a SMIL animation element ends.
//
// Category: SVG
//
// https://developer.mozilla.org/docs/Web/Events/endEvent
func OnEndEvent(listener func(e *TimeEvent)) markup.EventListener {
	return mkup.NewEventListener("endEvent", func(le *event.Event) {
		dispatchEventHelperTimeEvent(le, listener)
	})
}

// OnRepeatEvent is an event fired when a SMIL animation element is repeated.
//
// Category: SVG
//
// https://developer.mozilla.org/docs/Web/Events/repeatEvent
func OnRepeatEvent(listener func(e *TimeEvent)) markup.EventListener {
	return mkup.NewEventListener("repeatEvent", func(le *event.Event) {
		dispatchEventHelperTimeEvent(le, listener)
	})
}

// OnSVGAbort is an event fired when page loading has been stopped before the
// SVG was loaded.
//
// Category: SVG
//
// https://developer.mozilla.org/docs/Web/Events/SVGAbort
func OnSVGAbort(listener func(e *SVGEvent)) markup.EventListener {
	return mkup.NewEventListener("SVGAbort", func(le *event.Event) {
		dispatchEventHelperSVGEvent(le, listener)
	})
}

// OnSVGError is an event fired when an error has occurred before the SVG was
// loaded.
//
// Category: SVG
//
// https://developer.mozilla.org/docs/Web/Events/SVGError
func OnSVGError(listener func(e *SVGEvent)) markup.EventListener {
	return mkup.NewEventListener("SVGError", func(le *event.Event) {
		dispatchEventHelperSVGEvent(le, listener)
	})
}

// OnSVGLoad is an event fired when an SVG document has been loaded and parsed.
//
// Category: SVG
//
// https://developer.mozilla.org/docs/Web/Events/SVGLoad
func OnSVGLoad(listener func(e *SVGEvent)) markup.EventListener {
	return mkup.NewEventListener("SVGLoad", func(le *event.Event) {
		dispatchEventHelperSVGEvent(le, listener)
	})
}

// OnSVGResize is an event fired when an SVG document is being resized.
//
// Category: SVG
//
// https://developer.mozilla.org/docs/Web/Events/SVGResize
func OnSVGResize(listener func(e *SVGEvent)) markup.EventListener {
	return mkup.NewEventListener("SVGResize", func(le *event.Event) {
		dispatchEventHelperSVGEvent(le, listener)
	})
}

// OnSVGScroll is an event fired when an SVG document is being scrolled.
//
// Category: SVG
//
// https://developer.mozilla.org/docs/Web/Events/SVGScroll
func OnSVGScroll(listener func(e *SVGEvent)) markup.EventListener {
	return mkup.NewEventListener("SVGScroll", func(le *event.Event) {
		dispatchEventHelperSVGEvent(le, listener)
	})
}

// OnSVGUnload is an event fired when an SVG document has been removed from a
// window or frame.
//
// Category: SVG
//
// https://developer.mozilla.org/docs/Web/Events/SVGUnload
func OnSVGUnload(listener func(e *SVGEvent)) markup.EventListener {
	return mkup.NewEventListener("SVGUnload", func(le *event.Event) {
		dispatchEventHelperSVGEvent(le, listener)
	})
}

// OnSVGZoom is an event fired when an SVG document is being zoomed.
//
// Category: SVG
//
// https://developer.mozilla.org/docs/Web/Events/SVGZoom
func OnSVGZoom(listener func(e *SVGZoomEvent)) markup.EventListener {
	return mkup.NewEventListener("SVGZoom", func(le *event.Event) {
		dispatchEventHelperSVGZoomEvent(le, listener)
	})
}
