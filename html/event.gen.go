//go:generate go run gen_event.go

// Package event defines markup to bind DOM events.
//
// Generated from "Event reference" by Mozilla Contributors,
// https://developer.mozilla.org/en-US/docs/Web/Events, licensed under
// CC-BY-SA 2.5.
package html

import "github.com/yossoy/exciton/markup"
import "github.com/yossoy/exciton/event"

// AfterPrint is an event fired when the associated document has started
// printing or the print preview has been closed.
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/afterprint
func AfterPrint(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "afterprint", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// BeforePrint is an event fired when the associated document is about to be
// printed or previewed for printing.
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/beforeprint
func BeforePrint(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "beforeprint", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// BeforeUnload is an event fired when the window, the document and its
// resources are about to be unloaded.
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/beforeunload
func BeforeUnload(listener func(e *BeforeUnloadEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "beforeunload", Listener: func(le *event.Event) {
		dispatchEventHelperBeforeUnloadEvent(le, listener)
	}}
}

// CanPlay is an event fired when the user agent can play the media, but
// estimates that not enough data has been loaded to play the media up to its
// end without having to stop for further buffering of content.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/canplay
func CanPlay(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "canplay", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// CanPlayThrough is an event fired when the user agent can play the media up
// to its end without having to stop for further buffering of content.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/canplaythrough
func CanPlayThrough(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "canplaythrough", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// CompositionEnd is an event fired when the composition of a passage of text
// has been completed or canceled.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/compositionend
func CompositionEnd(listener func(e *CompositionEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "compositionend", Listener: func(le *event.Event) {
		dispatchEventHelperCompositionEvent(le, listener)
	}}
}

// CompositionStart is an event fired when the composition of a passage of text
// is prepared (similar to keydown for a keyboard input, but works with other
// inputs such as speech recognition).
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/compositionstart
func CompositionStart(listener func(e *CompositionEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "compositionstart", Listener: func(le *event.Event) {
		dispatchEventHelperCompositionEvent(le, listener)
	}}
}

// CompositionUpdate is an event fired when a character is added to a passage
// of text being composed.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/compositionupdate
func CompositionUpdate(listener func(e *CompositionEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "compositionupdate", Listener: func(le *event.Event) {
		dispatchEventHelperCompositionEvent(le, listener)
	}}
}

// ContextMenu is an event fired when the right button of the mouse is clicked
// (before the context menu is displayed).
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/contextmenu
func ContextMenu(listener func(e *MouseEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "contextmenu", Listener: func(le *event.Event) {
		dispatchEventHelperMouseEvent(le, listener)
	}}
}

// DoubleClick is an event fired when a pointing device button is clicked twice
// on an element.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/dblclick
func DoubleClick(listener func(e *MouseEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "dblclick", Listener: func(le *event.Event) {
		dispatchEventHelperMouseEvent(le, listener)
	}}
}

// DragEnd is an event fired when a drag operation is being ended (by releasing
// a mouse button or hitting the escape key).
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/dragend
func DragEnd(listener func(e *DragEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "dragend", Listener: func(le *event.Event) {
		dispatchEventHelperDragEvent(le, listener)
	}}
}

// DragEnter is an event fired when a dragged element or text selection enters
// a valid drop target.
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/dragenter
func DragEnter(listener func(e *DragEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "dragenter", Listener: func(le *event.Event) {
		dispatchEventHelperDragEvent(le, listener)
	}}
}

// DragLeave is an event fired when a dragged element or text selection leaves
// a valid drop target.
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/dragleave
func DragLeave(listener func(e *DragEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "dragleave", Listener: func(le *event.Event) {
		dispatchEventHelperDragEvent(le, listener)
	}}
}

// DragOver is an event fired when an element or text selection is being
// dragged over a valid drop target (every 350ms).
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/dragover
func DragOver(listener func(e *DragEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "dragover", Listener: func(le *event.Event) {
		dispatchEventHelperDragEvent(le, listener)
	}}
}

// DragStart is an event fired when the user starts dragging an element or text
// selection.
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/dragstart
func DragStart(listener func(e *DragEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "dragstart", Listener: func(le *event.Event) {
		dispatchEventHelperDragEvent(le, listener)
	}}
}

// DurationChange is an event fired when the duration attribute has been
// updated.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/durationchange
func DurationChange(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "durationchange", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// FocusIn is an event fired when an element is about to receive focus
// (bubbles).
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/focusin
func FocusIn(listener func(e *FocusEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "focusin", Listener: func(le *event.Event) {
		dispatchEventHelperFocusEvent(le, listener)
	}}
}

// FocusOut is an event fired when an element is about to lose focus (bubbles).
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/focusout
func FocusOut(listener func(e *FocusEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "focusout", Listener: func(le *event.Event) {
		dispatchEventHelperFocusEvent(le, listener)
	}}
}

// HashChange is an event fired when the fragment identifier of the URL has
// changed (the part of the URL after the #).
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/hashchange
func HashChange(listener func(e *HashChangeEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "hashchange", Listener: func(le *event.Event) {
		dispatchEventHelperHashChangeEvent(le, listener)
	}}
}

// KeyDown is an event fired when a key is pressed down.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/keydown
func KeyDown(listener func(e *KeyboardEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "keydown", Listener: func(le *event.Event) {
		dispatchEventHelperKeyboardEvent(le, listener)
	}}
}

// KeyPress is an event fired when a key is pressed down and that key normally
// produces a character value (use input instead).
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/keypress
func KeyPress(listener func(e *KeyboardEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "keypress", Listener: func(le *event.Event) {
		dispatchEventHelperKeyboardEvent(le, listener)
	}}
}

// KeyUp is an event fired when a key is released.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/keyup
func KeyUp(listener func(e *KeyboardEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "keyup", Listener: func(le *event.Event) {
		dispatchEventHelperKeyboardEvent(le, listener)
	}}
}

// LanguageChange is an event fired when the user's preferred languages have
// changed.
//
// Category: HTML 5.1The definition of 'NavigatorLanguage.languages' in that specification.
//
// https://developer.mozilla.org/docs/Web/Events/languagechange
func LanguageChange(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "languagechange", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// LoadedData is an event fired when the first frame of the media has finished
// loading.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/loadeddata
func LoadedData(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "loadeddata", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// LoadedMetadata is an event fired when the metadata has been loaded.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/loadedmetadata
func LoadedMetadata(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "loadedmetadata", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// MouseDown is an event fired when a pointing device button (usually a mouse)
// is pressed on an element.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/mousedown
func MouseDown(listener func(e *MouseEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "mousedown", Listener: func(le *event.Event) {
		dispatchEventHelperMouseEvent(le, listener)
	}}
}

// MouseEnter is an event fired when a pointing device is moved onto the
// element that has the listener attached.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/mouseenter
func MouseEnter(listener func(e *MouseEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "mouseenter", Listener: func(le *event.Event) {
		dispatchEventHelperMouseEvent(le, listener)
	}}
}

// MouseLeave is an event fired when a pointing device is moved off the element
// that has the listener attached.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/mouseleave
func MouseLeave(listener func(e *MouseEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "mouseleave", Listener: func(le *event.Event) {
		dispatchEventHelperMouseEvent(le, listener)
	}}
}

// MouseMove is an event fired when a pointing device is moved over an element.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/mousemove
func MouseMove(listener func(e *MouseEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "mousemove", Listener: func(le *event.Event) {
		dispatchEventHelperMouseEvent(le, listener)
	}}
}

// MouseOut is an event fired when a pointing device is moved off the element
// that has the listener attached or off one of its children.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/mouseout
func MouseOut(listener func(e *MouseEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "mouseout", Listener: func(le *event.Event) {
		dispatchEventHelperMouseEvent(le, listener)
	}}
}

// MouseOver is an event fired when a pointing device is moved onto the element
// that has the listener attached or onto one of its children.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/mouseover
func MouseOver(listener func(e *MouseEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "mouseover", Listener: func(le *event.Event) {
		dispatchEventHelperMouseEvent(le, listener)
	}}
}

// MouseUp is an event fired when a pointing device button is released over an
// element.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/mouseup
func MouseUp(listener func(e *MouseEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "mouseup", Listener: func(le *event.Event) {
		dispatchEventHelperMouseEvent(le, listener)
	}}
}

// NoUpdate is an event fired when the manifest hadn't changed.
//
// Category: Offline
//
// https://developer.mozilla.org/docs/Web/Events/noupdate
func NoUpdate(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "noupdate", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnAbort is an event fired when the loading of a resource has been aborted.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/abort
func OnAbort(listener func(e *UIEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "abort", Listener: func(le *event.Event) {
		dispatchEventHelperUIEvent(le, listener)
	}}
}

// OnBlur is an event fired when an element has lost focus (does not bubble).
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/blur
func OnBlur(listener func(e *FocusEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "blur", Listener: func(le *event.Event) {
		dispatchEventHelperFocusEvent(le, listener)
	}}
}

// OnCached is an event fired when the resources listed in the manifest have
// been downloaded, and the application is now cached.
//
// Category: Offline
//
// https://developer.mozilla.org/docs/Web/Events/cached
func OnCached(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "cached", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnChange is an event fired when the change event is fired for <input>,
// <select>, and <textarea> elements when a change to the element's value is
// committed by the user.
//
// Category: DOM L2, HTML5
//
// https://developer.mozilla.org/docs/Web/Events/change
func OnChange(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "change", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnChecking is an event fired when the user agent is checking for an update,
// or attempting to download the cache manifest for the first time.
//
// Category: Offline
//
// https://developer.mozilla.org/docs/Web/Events/checking
func OnChecking(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "checking", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnClick is an event fired when a pointing device button has been pressed and
// released on an element.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/click
func OnClick(listener func(e *MouseEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "click", Listener: func(le *event.Event) {
		dispatchEventHelperMouseEvent(le, listener)
	}}
}

// OnDOMContentLoaded is an event fired when the document has finished loading
// (but not its dependent resources).
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/DOMContentLoaded
func OnDOMContentLoaded(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "DOMContentLoaded", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnDownloading is an event fired when the user agent has found an update and
// is fetching it, or is downloading the resources listed by the cache manifest
// for the first time.
//
// Category: Offline
//
// https://developer.mozilla.org/docs/Web/Events/downloading
func OnDownloading(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "downloading", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnDrag is an event fired when an element or text selection is being dragged
// (every 350ms).
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/drag
func OnDrag(listener func(e *DragEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "drag", Listener: func(le *event.Event) {
		dispatchEventHelperDragEvent(le, listener)
	}}
}

// OnDrop is an event fired when an element is dropped on a valid drop target.
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/drop
func OnDrop(listener func(e *DragEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "drop", Listener: func(le *event.Event) {
		dispatchEventHelperDragEvent(le, listener)
	}}
}

// OnEmptied is an event fired when the media has become empty; for example,
// this event is sent if the media has already been loaded (or partially
// loaded), and the load() method is called to reload it.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/emptied
func OnEmptied(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "emptied", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnEnded is an event fired when playback has stopped because the end of the
// media was reached.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/ended
func OnEnded(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "ended", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnError is an event fired when an error occurred while downloading the cache
// manifest or updating the content of the application.
//
// Category: Offline
//
// https://developer.mozilla.org/docs/Web/Events/error
func OnError(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "error", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnFocus is an event fired when an element has received focus (does not
// bubble).
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/focus
func OnFocus(listener func(e *FocusEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "focus", Listener: func(le *event.Event) {
		dispatchEventHelperFocusEvent(le, listener)
	}}
}

// OnInput is an event fired when the value of an element changes or the
// content of an element with the attribute contenteditable is modified.
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/input
func OnInput(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "input", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnInvalid is an event fired when a submittable element has been checked and
// doesn't satisfy its constraints.
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/invalid
func OnInvalid(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "invalid", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnLoad is an event fired when a resource and its dependent resources have
// finished loading.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/load
func OnLoad(listener func(e *UIEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "load", Listener: func(le *event.Event) {
		dispatchEventHelperUIEvent(le, listener)
	}}
}

// OnObsolete is an event fired when the manifest was found to have become a
// 404 or 410 page, so the application cache is being deleted.
//
// Category: Offline
//
// https://developer.mozilla.org/docs/Web/Events/obsolete
func OnObsolete(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "obsolete", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnOffline is an event fired when the browser has lost access to the network.
//
// Category: HTML5 offline
//
// https://developer.mozilla.org/docs/Web/Events/offline
func OnOffline(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "offline", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnOnline is an event fired when the browser has gained access to the network
// (but particular websites might be unreachable).
//
// Category: HTML5 offline
//
// https://developer.mozilla.org/docs/Web/Events/online
func OnOnline(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "online", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnPause is an event fired when playback has been paused.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/pause
func OnPause(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "pause", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnPlay is an event fired when playback has begun.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/play
func OnPlay(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "play", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnPlaying is an event fired when playback is ready to start after having
// been paused or delayed due to lack of data.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/playing
func OnPlaying(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "playing", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnProgress is an event fired when the user agent is downloading resources
// listed by the manifest.
//
// Category: Offline
//
// https://developer.mozilla.org/docs/Web/Reference/Events/progress_(appcache_event)
func OnProgress(listener func(e *ProgressEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "progress", Listener: func(le *event.Event) {
		dispatchEventHelperProgressEvent(le, listener)
	}}
}

// OnReset is an event fired when a form is reset.
//
// Category: DOM L2, HTML5
//
// https://developer.mozilla.org/docs/Web/Events/reset
func OnReset(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "reset", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnResize is an event fired when the document view has been resized.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/resize
func OnResize(listener func(e *UIEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "resize", Listener: func(le *event.Event) {
		dispatchEventHelperUIEvent(le, listener)
	}}
}

// OnScroll is an event fired when the document view or an element has been
// scrolled.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/scroll
func OnScroll(listener func(e *UIEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "scroll", Listener: func(le *event.Event) {
		dispatchEventHelperUIEvent(le, listener)
	}}
}

// OnSeeked is an event fired when a seek operation completed.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/seeked
func OnSeeked(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "seeked", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnSeeking is an event fired when a seek operation began.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/seeking
func OnSeeking(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "seeking", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnSelect is an event fired when some text is being selected.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/select
func OnSelect(listener func(e *UIEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "select", Listener: func(le *event.Event) {
		dispatchEventHelperUIEvent(le, listener)
	}}
}

// OnShow is an event fired when a contextmenu event was fired on/bubbled to an
// element that has a contextmenu attribute
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/show
func OnShow(listener func(e *MouseEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "show", Listener: func(le *event.Event) {
		dispatchEventHelperMouseEvent(le, listener)
	}}
}

// OnStalled is an event fired when the user agent is trying to fetch media
// data, but data is unexpectedly not forthcoming.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/stalled
func OnStalled(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "stalled", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnSubmit is an event fired when a form is submitted.
//
// Category: DOM L2, HTML5
//
// https://developer.mozilla.org/docs/Web/Events/submit
func OnSubmit(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "submit", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnSuspend is an event fired when media data loading has been suspended.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/suspend
func OnSuspend(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "suspend", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnUnload is an event fired when the document or a dependent resource is
// being unloaded.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/unload
func OnUnload(listener func(e *UIEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "unload", Listener: func(le *event.Event) {
		dispatchEventHelperUIEvent(le, listener)
	}}
}

// OnWaiting is an event fired when playback has stopped because of a temporary
// lack of data.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/waiting
func OnWaiting(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "waiting", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// OnWheel is an event fired when a wheel button of a pointing device is
// rotated in any direction.
//
// Category: DOM L3
//
// https://developer.mozilla.org/docs/Web/Events/wheel
func OnWheel(listener func(e *WheelEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "wheel", Listener: func(le *event.Event) {
		dispatchEventHelperWheelEvent(le, listener)
	}}
}

// PageHide is an event fired when a session history entry is being traversed
// from.
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/pagehide
func PageHide(listener func(e *PageTransitionEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "pagehide", Listener: func(le *event.Event) {
		dispatchEventHelperPageTransitionEvent(le, listener)
	}}
}

// PageShow is an event fired when a session history entry is being traversed
// to.
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/pageshow
func PageShow(listener func(e *PageTransitionEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "pageshow", Listener: func(le *event.Event) {
		dispatchEventHelperPageTransitionEvent(le, listener)
	}}
}

// PopState is an event fired when a session history entry is being navigated
// to (in certain cases).
//
// Category: HTML5
//
// https://developer.mozilla.org/docs/Web/Events/popstate
func PopState(listener func(e *PopStateEvent)) *markup.EventListener {
	return &markup.EventListener{Name: "popstate", Listener: func(le *event.Event) {
		dispatchEventHelperPopStateEvent(le, listener)
	}}
}

// RateChange is an event fired when the playback rate has changed.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/ratechange
func RateChange(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "ratechange", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// ReadyStateChange is an event fired when the readyState attribute of a
// document has changed.
//
// Category: HTML5 and XMLHttpRequest
//
// https://developer.mozilla.org/docs/Web/Events/readystatechange
func ReadyStateChange(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "readystatechange", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// SelectStart is an event fired when a selection just started.
//
// Category: Selection API
//
// https://developer.mozilla.org/docs/Web/Events/selectstart
func SelectStart(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "selectstart", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// SelectionChange is an event fired when the selection in the document has
// been changed.
//
// Category: Selection API
//
// https://developer.mozilla.org/docs/Web/Events/selectionchange
func SelectionChange(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "selectionchange", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// TimeUpdate is an event fired when the time indicated by the currentTime
// attribute has been updated.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/timeupdate
func TimeUpdate(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "timeupdate", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// UpdateReady is an event fired when the resources listed in the manifest have
// been newly redownloaded, and the script can use swapCache() to switch to the
// new cache.
//
// Category: Offline
//
// https://developer.mozilla.org/docs/Web/Events/updateready
func UpdateReady(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "updateready", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}

// VolumeChange is an event fired when the volume has changed.
//
// Category: HTML5 media
//
// https://developer.mozilla.org/docs/Web/Events/volumechange
func VolumeChange(listener func(e *Event)) *markup.EventListener {
	return &markup.EventListener{Name: "volumechange", Listener: func(le *event.Event) {
		dispatchEventHelperEvent(le, listener)
	}}
}
