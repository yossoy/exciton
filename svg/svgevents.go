package svg

import (
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/html"
	"github.com/yossoy/exciton/markup"
)

type SVGEvent struct {
	html.Event
}

func dispatchEventHelperSVGEvent(ee *event.Event, handler func(e *SVGEvent)) {
	var e SVGEvent
	if err := ee.Argument.Decode(&e); err != nil {
		panic(err)
	}
	handler(&e)
}

type SVGPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type SVGRect struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type SVGZoomEvent struct {
	html.UIEvent

	ZoomRectScreen    SVGRect  `json:"zoomRectScreen"`
	PreviousScale     float64  `json:"previousScale"`
	PreviousTranslate SVGPoint `json:"previousTranslate"`
	NewScale          float64  `json:"newScale"`
	NewTranslate      SVGPoint `json:"newTranslate"`
}

func dispatchEventHelperSVGZoomEvent(ee *event.Event, handler func(e *SVGZoomEvent)) {
	var e SVGZoomEvent
	if err := ee.Argument.Decode(&e); err != nil {
		panic(err)
	}
	handler(&e)
}

type TimeEvent struct {
	html.Event

	Detail int64               `json:"detail"`         // Specifies some detail information about the Event, depending on the type of the event. For this event type, indicates the repeat number for the animation.
	View   *markup.EventTarget `json:"view,omitempty"` // The view attribute identifies the AbstractView [DOM2VIEWS] from which the event was generated.
}

func dispatchEventHelperTimeEvent(ee *event.Event, handler func(e *TimeEvent)) {
	var e TimeEvent
	if err := ee.Argument.Decode(&e); err != nil {
		panic(err)
	}
	handler(&e)
}
