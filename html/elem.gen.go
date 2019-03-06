//go:generate go run gen_elem.go

// Package html defines markup to create DOM elements.
//
// Generated from "HTML element reference" by Mozilla Contributors,
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element, licensed under
// CC-BY-SA 2.5.
package html

import mkup "github.com/yossoy/exciton/markup"

// Anchor (or anchor element) creates a hyperlink to other web pages, files,
// locations within the same page, email addresses, or any other URL.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/a
func Anchor(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("a", markup...)
}

// The HTML Abbreviation element (<abbr>) represents an abbreviation or
// acronym; the optional title attribute can provide an expansion or
// description for the abbreviation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/abbr
func Abbreviation(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("abbr", markup...)
}

// Address indicates that the enclosed HTML provides contact information for a
// person or people, or for an organization.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/address
func Address(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("address", markup...)
}

// Area defines a hot-spot region on an image, and optionally associates it
// with a hypertext link. This element is used only within a <map> element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/area
func Area(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("area", markup...)
}

// Article represents a self-contained composition in a document, page,
// application, or site, which is intended to be independently distributable or
// reusable (e.g., in syndication). Examples include: a forum post, a magazine
// or newspaper article, or a blog entry.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/article
func Article(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("article", markup...)
}

// Aside represents a portion of a document whose content is only indirectly
// related to the document's main content.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/aside
func Aside(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("aside", markup...)
}

// Audio is used to embed sound content in documents. It may contain one or
// more audio sources, represented using the src attribute or the <source>
// element: the browser will choose the most suitable one. It can also be the
// destination for streamed media, using a MediaStream.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/audio
func Audio(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("audio", markup...)
}

// The HTML Bring Attention To element (<b>) is used to draw the reader's
// attention to the element's contents, which are not otherwise granted special
// importance.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/b
func Bold(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("b", markup...)
}

// Base specifies the base URL to use for all relative URLs contained within a
// document. There can be only one <base> element in a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/base
func Base(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("base", markup...)
}

// The HTML Bidirectional Isolate element (<bdi>) tells the browser's
// bidirectional algorithm to treat the text it contains in isolation from its
// surrounding text.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/bdi
func BidirectionalIsolation(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("bdi", markup...)
}

// The HTML Bidirectional Text Override element (<bdo>) overrides the current
// directionality of text, so that the text within is rendered in a different
// direction.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/bdo
func BidirectionalOverride(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("bdo", markup...)
}

// BlockQuote (or HTML Block Quotation Element) indicates that the enclosed
// text is an extended quotation. Usually, this is rendered visually by
// indentation (see Notes for how to change it). A URL for the source of the
// quotation may be given using the cite attribute, while a text representation
// of the source can be given using the <cite> element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/blockquote
func BlockQuote(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("blockquote", markup...)
}

// Body represents the content of an HTML document. There can be only one
// <body> element in a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/body
func Body(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("body", markup...)
}

// Break produces a line break in text (carriage-return). It is useful for
// writing a poem or an address, where the division of lines is significant.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/br
func Break(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("br", markup...)
}

// Button represents a clickable button, which can be used in forms or anywhere
// in a document that needs simple, standard button functionality.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/button
func Button(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("button", markup...)
}

// Use the HTML <canvas> element with either the canvas scripting API or the
// WebGL API to draw graphics and animations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/canvas
func Canvas(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("canvas", markup...)
}

// The HTML Table Caption element (<caption>) specifies the caption (or title)
// of a table, and if used is always the first child of a <table>.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/caption
func Caption(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("caption", markup...)
}

// The HTML Citation element (<cite>) is used to describe a reference to a
// cited creative work, and must include either the title or author or the URL
// of that work.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/cite
func Citation(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("cite", markup...)
}

// Code displays its contents styled in a fashion intended to indicate that the
// text is a short fragment of computer code.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/code
func Code(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("code", markup...)
}

// Column defines a column within a table and is used for defining common
// semantics on all common cells. It is generally found within a <colgroup>
// element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/col
func Column(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("col", markup...)
}

// ColumnGroup defines a group of columns within a table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/colgroup
func ColumnGroup(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("colgroup", markup...)
}

// Data links a given content with a machine-readable translation. If the
// content is time- or date-related, the <time> element must be used.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/data
func Data(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("data", markup...)
}

// DataList contains a set of <option> elements that represent the values
// available for other controls.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/datalist
func DataList(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("datalist", markup...)
}

// Description provides the details about or the definition of the preceding
// term (<dt>) in a description list (<dl>).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dd
func Description(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("dd", markup...)
}

// DeletedText represents a range of text that has been deleted from a
// document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/del
func DeletedText(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("del", markup...)
}

// The HTML Details Element (<details>) creates a disclosure widget in which
// information is visible only when the widget is toggled into an "open" state.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/details
func Details(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("details", markup...)
}

// The HTML Definition element (<dfn>) is used to indicate the term being
// defined within the context of a definition phrase or sentence.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dfn
func Definition(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("dfn", markup...)
}

// Dialog represents a dialog box or other interactive component, such as an
// inspector or window.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dialog
func Dialog(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("dialog", markup...)
}

// The HTML Content Division element (<div>) is the generic container for flow
// content. It has no effect on the content or layout until styled using CSS.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/div
func Div(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("div", markup...)
}

// DescriptionList represents a description list. The element encloses a list
// of groups of terms (specified using the <dt> element) and descriptions
// (provided by <dd> elements). Common uses for this element are to implement a
// glossary or to display metadata (a list of key-value pairs).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dl
func DescriptionList(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("dl", markup...)
}

// DefinitionTerm specifies a term in a description or definition list, and as
// such must be used inside a <dl> element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dt
func DefinitionTerm(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("dt", markup...)
}

// Emphasis marks text that has stress emphasis. The <em> element can be
// nested, with each level of nesting indicating a greater degree of emphasis.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/em
func Emphasis(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("em", markup...)
}

// Embed embeds external content at the specified point in the document. This
// content is provided by an external application or other source of
// interactive content such as a browser plug-in.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/embed
func Embed(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("embed", markup...)
}

// FieldSet is used to group several controls as well as labels (<label>)
// within a web form.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/fieldset
func FieldSet(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("fieldset", markup...)
}

// The HTML <figcaption> or Figure Caption element represents a caption or
// legend describing the rest of the contents of its parent <figure> element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/figcaption
func FigureCaption(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("figcaption", markup...)
}

// The HTML <figure> (Figure With Optional Caption) element represents
// self-contained content, potentially with an optional caption, which is
// specified using the (<figcaption>) element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/figure
func Figure(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("figure", markup...)
}

// Footer represents a footer for its nearest sectioning content or sectioning
// root element. A footer typically contains information about the author of
// the section, copyright data or links to related documents.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/footer
func Footer(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("footer", markup...)
}

// Form represents a document section that contains interactive controls for
// submitting information to a web server.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/form
func Form(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("form", markup...)
}

// The HTML <h1>–<h6> elements represent six levels of section headings. <h1>
// is the highest section level and <h6> is the lowest.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/Heading_Elements
func Heading1(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("h1", markup...)
}

// Header represents introductory content, typically a group of introductory or
// navigational aids. It may contain some heading elements but also a logo, a
// search form, an author name, and other elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/header
func Header(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("header", markup...)
}

// HeadingsGroup represents a multi-level heading for a section of a document.
// It groups a set of <h1>–<h6> elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/hgroup
func HeadingsGroup(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("hgroup", markup...)
}

// HorizontalRule represents a thematic break between paragraph-level elements:
// for example, a change of scene in a story, or a shift of topic within a
// section.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/hr
func HorizontalRule(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("hr", markup...)
}

// Italic represents a range of text that is set off from the normal text for
// some reason. Some examples include technical terms, foreign language
// phrases, or fictional character thoughts. It is typically displayed in
// italic type.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/i
func Italic(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("i", markup...)
}

// The HTML Inline Frame element (<iframe>) represents a nested browsing
// context, embedding another HTML page into the current one.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/iframe
func InlineFrame(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("iframe", markup...)
}

// Image embeds an image into the document. It is a replaced element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/img
func Image(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("img", markup...)
}

// Input is used to create interactive controls for web-based forms in order to
// accept data from the user; a wide variety of types of input data and control
// widgets are available, depending on the device and user agent.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input
func Input(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("input", markup...)
}

// InsertedText represents a range of text that has been added to a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/ins
func InsertedText(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("ins", markup...)
}

// The HTML Keyboard Input element (<kbd>) represents a span of inline text
// denoting textual user input from a keyboard, voice input, or any other text
// entry device.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/kbd
func KeyboardInput(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("kbd", markup...)
}

// Label represents a caption for an item in a user interface.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/label
func Label(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("label", markup...)
}

// Legend represents a caption for the content of its parent <fieldset>.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/legend
func Legend(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("legend", markup...)
}

// ListItem is used to represent an item in a list. It must be contained in a
// parent element: an ordered list (<ol>), an unordered list (<ul>), or a menu
// (<menu>). In menus and unordered lists, list items are usually displayed
// using bullet points. In ordered lists, they are usually displayed with an
// ascending counter on the left, such as a number or letter.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/li
func ListItem(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("li", markup...)
}

// The HTML External Resource Link element (<link>) specifies relationships
// between the current document and an external resource. This element is most
// commonly used to link to stylesheets, but is also used to establish site
// icons (both "favicon" style icons and mobile home screen/app icons) among
// other things.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/link
func Link(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("link", markup...)
}

// Main represents the dominant content of the <body> of a document. The main
// content area consists of content that is directly related to or expands upon
// the central topic of a document, or the central functionality of an
// application.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/main
func Main(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("main", markup...)
}

// Map is used with <area> elements to define an image map (a clickable link
// area).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/map
func Map(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("map", markup...)
}

// The HTML Mark Text element (<mark>) represents text which is marked or
// highlighted for reference or notation purposes, due to the marked passage's
// relevance or importance in the enclosing context.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/mark
func Mark(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("mark", markup...)
}

// Menu represents a group of commands that a user can perform or activate.
// This includes both list menus, which might appear across the top of a
// screen, as well as context menus, such as those that might appear underneath
// a button after it has been clicked.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/menu
func Menu(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("menu", markup...)
}

// Meta represents metadata that cannot be represented by other HTML
// meta-related elements, like <base>, <link>, <script>, <style> or <title>.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta
func Meta(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("meta", markup...)
}

// Meter represents either a scalar value within a known range or a fractional
// value.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meter
func Meter(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("meter", markup...)
}

// Navigation represents a section of a page whose purpose is to provide
// navigation links, either within the current document or to other documents.
// Common examples of navigation sections are menus, tables of contents, and
// indexes.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/nav
func Navigation(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("nav", markup...)
}

// NoScript defines a section of HTML to be inserted if a script type on the
// page is unsupported or if scripting is currently turned off in the browser.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/noscript
func NoScript(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("noscript", markup...)
}

// Object represents an external resource, which can be treated as an image, a
// nested browsing context, or a resource to be handled by a plugin.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/object
func Object(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("object", markup...)
}

// OrderedList represents an ordered list of items, typically rendered as a
// numbered list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/ol
func OrderedList(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("ol", markup...)
}

// OptionsGroup creates a grouping of options within a <select> element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/optgroup
func OptionsGroup(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("optgroup", markup...)
}

// Option is used to define an item contained in a <select>, an <optgroup>, or
// a <datalist> element. As such, <option> can represent menu items in popups
// and other lists of items in an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/option
func Option(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("option", markup...)
}

// The HTML Output element (<output>) is a container element into which a site
// or app can inject the results of a calculation or the outcome of a user
// action.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/output
func Output(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("output", markup...)
}

// Paragraph represents a paragraph.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/p
func Paragraph(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("p", markup...)
}

// Parameter defines parameters for an <object> element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/param
func Parameter(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("param", markup...)
}

// Picture contains zero or more <source> elements and one <img> element to
// provide versions of an image for different display/device scenarios.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/picture
func Picture(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("picture", markup...)
}

// Preformatted represents preformatted text which is to be presented exactly
// as written in the HTML file.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/pre
func Preformatted(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("pre", markup...)
}

// Progress displays an indicator showing the completion progress of a task,
// typically displayed as a progress bar.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/progress
func Progress(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("progress", markup...)
}

// Quote indicates that the enclosed text is a short inline quotation. Most
// modern browsers implement this by surrounding the text in quotation marks.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/q
func Quote(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("q", markup...)
}

// The HTML Ruby Base (<rb>) element is used to delimit the base text component
// of a <ruby> annotation, i.e. the text that is being annotated.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/rb
func Rb(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("rb", markup...)
}

// The HTML Ruby Fallback Parenthesis (<rp>) element is used to provide
// fall-back parentheses for browsers that do not support display of ruby
// annotations using the <ruby> element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/rp
func RubyParenthesis(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("rp", markup...)
}

// The HTML Ruby Text (<rt>) element specifies the ruby text component of a
// ruby annotation, which is used to provide pronunciation, translation, or
// transliteration information for East Asian typography. The <rt> element must
// always be contained within a <ruby> element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/rt
func RubyText(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("rt", markup...)
}

// The HTML Ruby Text Container (<rtc>) element embraces semantic annotations
// of characters presented in a ruby of <rb> elements used inside of <ruby>
// element. <rb> elements can have both pronunciation (<rt>) and semantic
// (<rtc>) annotations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/rtc
func RubyTextContainer(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("rtc", markup...)
}

// Ruby represents a ruby annotation. Ruby annotations are for showing
// pronunciation of East Asian characters.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/ruby
func Ruby(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("ruby", markup...)
}

// Strikethrough renders text with a strikethrough, or a line through it. Use
// the <s> element to represent things that are no longer relevant or no longer
// accurate. However, <s> is not appropriate when indicating document edits;
// for that, use the <del> and <ins> elements, as appropriate.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/s
func Strikethrough(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("s", markup...)
}

// The HTML Sample Element (<samp>) is used to enclose inline text which
// represents sample (or quoted) output from a computer program.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/samp
func Sample(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("samp", markup...)
}

// Script is used to embed or reference executable code; this is typically used
// to embed or refer to JavaScript code.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/script
func Script(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("script", markup...)
}

// Section represents a standalone section — which doesn't have a more
// specific semantic element to represent it — contained within an HTML
// document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/section
func Section(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("section", markup...)
}

// Select represents a control that provides a menu of options
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/select
func Select(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("select", markup...)
}

// Slot—part of the Web Components technology suite—is a placeholder inside
// a web component that you can fill with your own markup, which lets you
// create separate DOM trees and present them together.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/slot
func Slot(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("slot", markup...)
}

// Small makes the text font size one size smaller (for example, from large to
// medium, or from small to x-small) down to the browser's minimum font size.
// In HTML5, this element is repurposed to represent side-comments and small
// print, including copyright and legal text, independent of its styled
// presentation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/small
func Small(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("small", markup...)
}

// Source specifies multiple media resources for the <picture>, the <audio>
// element, or the <video> element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/source
func Source(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("source", markup...)
}

// Span is a generic inline container for phrasing content, which does not
// inherently represent anything. It can be used to group elements for styling
// purposes (using the class or id attributes), or because they share attribute
// values, such as lang.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/span
func Span(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("span", markup...)
}

// The HTML Strong Importance Element (<strong>) indicates that its contents
// have strong importance, seriousness, or urgency. Browsers typically render
// the contents in bold type.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/strong
func Strong(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("strong", markup...)
}

// Style contains style information for a document, or part of a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/style
func Style(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("style", markup...)
}

// The HTML Subscript element (<sub>) specifies inline text which should be
// displayed as subscript for solely typographical reasons.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/sub
func Subscript(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("sub", markup...)
}

// The HTML Disclosure Summary element (<summary>) element specifies a summary,
// caption, or legend for a <details> element's disclosure box.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/summary
func Summary(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("summary", markup...)
}

// The HTML Superscript element (<sup>) specifies inline text which is to be
// displayed as superscript for solely typographical reasons.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/sup
func Superscript(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("sup", markup...)
}

// Table represents tabular data — that is, information presented in a
// two-dimensional table comprised of rows and columns of cells containing
// data.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/table
func Table(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("table", markup...)
}

// The HTML Table Body element (<tbody>) encapsulates a set of table row (<tr>
// elements), indicating that they comprise the body of the table (<table>).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/tbody
func TableBody(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("tbody", markup...)
}

// TableData defines a cell of a table that contains data. It participates in
// the table model.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/td
func TableData(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("td", markup...)
}

// The HTML Content Template (<template>) element is a mechanism for holding
// client-side content that is not to be rendered when a page is loaded but may
// subsequently be instantiated during runtime using JavaScript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/template
func Template(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("template", markup...)
}

// TextArea represents a multi-line plain-text editing control, useful when you
// want to allow users to enter a sizeable amount of free-form text, for
// example a comment on a review or feedback form.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/textarea
func TextArea(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("textarea", markup...)
}

// TableFoot defines a set of rows summarizing the columns of the table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/tfoot
func TableFoot(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("tfoot", markup...)
}

// TableHeader defines a cell as header of a group of table cells. The exact
// nature of this group is defined by the scope and headers attributes.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/th
func TableHeader(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("th", markup...)
}

// TableHead defines a set of rows defining the head of the columns of the
// table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/thead
func TableHead(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("thead", markup...)
}

// Time represents a specific period in time.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/time
func Time(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("time", markup...)
}

// The HTML Title element (<title>) defines the document's title that is shown
// in a browser's title bar or a page's tab.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/title
func Title(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("title", markup...)
}

// TableRow defines a row of cells in a table. The row's cells can then be
// established using a mix of <td> (data cell) and <th> (header cell)
// elements.The HTML <tr> element specifies that the markup contained inside
// the <tr> block comprises one row of a table, inside which the <th> and <td>
// elements create header and data cells, respectively, within the row.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/tr
func TableRow(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("tr", markup...)
}

// Track is used as a child of the media elements <audio> and <video>. It lets
// you specify timed text tracks (or time-based data), for example to
// automatically handle subtitles. The tracks are formatted in WebVTT format
// (.vtt files) — Web Video Text Tracks or Timed Text Markup Language (TTML).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/track
func Track(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("track", markup...)
}

// The HTML Unarticulated Annotation Element (<u>) represents a span of inline
// text which should be rendered in a way that indicates that it has a
// non-textual annotation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/u
func Underline(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("u", markup...)
}

// UnorderedList represents an unordered list of items, typically rendered as a
// bulleted list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/ul
func UnorderedList(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("ul", markup...)
}

// The HTML Variable element (<var>) represents the name of a variable in a
// mathematical expression or a programming context.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/var
func Variable(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("var", markup...)
}

// The HTML Video element (<video>) embeds a media player which supports video
// playback into the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/video
func Video(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("video", markup...)
}

// WordBreakOpportunity represents a word break opportunity—a position within
// text where the browser may optionally break a line, though its line-breaking
// rules would not otherwise create a break at that location.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/wbr
func WordBreakOpportunity(markup ...mkup.MarkupOrChild) mkup.RenderResult {
	return mkup.Tag("wbr", markup...)
}
