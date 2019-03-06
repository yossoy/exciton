// Package svg defines markup to create SVG elements.
//
// Generated from "SVG element reference" by Mozilla Contributors,
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element, licensed under
// CC-BY-SA 2.5.
package svg

import "github.com/yossoy/exciton/markup"

// Anchor creates a hyperlink to other web pages, files, locations within the
// same page, email addresses, or any other URL.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/a
func Anchor(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("a", SVGNamespace, children...)
}

// AltGlyph allows sophisticated selection of the glyphs used to render its
// child character data.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/altGlyph
func AltGlyph(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("altGlyph", SVGNamespace, children...)
}

// AltGlyphDef defines a substitution representation for glyphs.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/altGlyphDef
func AltGlyphDef(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("altGlyphDef", SVGNamespace, children...)
}

// The <altGlyphItem> element provides a set of candidates for glyph
// substitution by the <altGlyph> element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/altGlyphItem
func AltGlyphItem(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("altGlyphItem", SVGNamespace, children...)
}

// This element implements the SVGAnimateElement interface.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/animate
func Animate(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("animate", SVGNamespace, children...)
}

// AnimateColor specifies a color transformation over time.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/animateColor
func AnimateColor(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("animateColor", SVGNamespace, children...)
}

// The <animateMotion> element causes a referenced element to move along a
// motion path.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/animateMotion
func AnimateMotion(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("animateMotion", SVGNamespace, children...)
}

// The animateTransform element animates a transformation attribute on a target
// element, thereby allowing animations to control translation, scaling,
// rotation and/or skewing.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/animateTransform
func AnimateTransform(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("animateTransform", SVGNamespace, children...)
}

// Circle is an SVG basic shape, used to create circles based on a center point
// and a radius.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/circle
func Circle(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("circle", SVGNamespace, children...)
}

// ClipPath defines a clipping path. A clipping path is used/referenced using
// the clip-path property.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/clipPath
func ClipPath(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("clipPath", SVGNamespace, children...)
}

// The <color-profile> element allows describing the color profile used for the
// image.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/color-profile
func ColorProfile(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("color-profile", SVGNamespace, children...)
}

// Cursor can be used to define a platform-independent custom cursor. A
// recommended approach for defining a platform-independent custom cursor is to
// create a PNG image and define a cursor element that references the PNG image
// and identifies the exact position within the image which is the pointer
// position (i.e., the hot spot).
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/cursor
func Cursor(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("cursor", SVGNamespace, children...)
}

// The <defs> element is used to store graphical objects that will be used at a
// later time. Objects created inside a <defs> element are not rendered
// directly. To display them you have to reference them (with a <use> element
// for example).
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/defs
func Defs(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("defs", SVGNamespace, children...)
}

// Each container element or graphics element in an SVG drawing can supply a
// description string using the <desc> element where the description is
// text-only.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/desc
func Desc(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("desc", SVGNamespace, children...)
}

// The <ellipse> element is an SVG basic shape, used to create ellipses based
// on a center coordinate, and both their x and y radius.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/ellipse
func Ellipse(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("ellipse", SVGNamespace, children...)
}

// The <feBlend> SVG filter primitive composes two objects together ruled by a
// certain blending mode. This is similar to what is known from image editing
// software when blending two layers. The mode is defined by the mode
// attribute.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feBlend
func FeBlend(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feBlend", SVGNamespace, children...)
}

// The <feColorMatrix> SVG filter element changes colors based on a
// transformation matrix. Every pixel's color value (represented by an
// [R,G,B,A] vector) is matrix multiplied to create a new color:
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feColorMatrix
func FeColorMatrix(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feColorMatrix", SVGNamespace, children...)
}

// Th <feComponentTransfer> SVG filter primitive performs color-component-wise
// remapping of data for each pixel. It allows operations like brightness
// adjustment, contrast adjustment, color balance or thresholding.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feComponentTransfer
func FeComponentTransfer(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feComponentTransfer", SVGNamespace, children...)
}

// The <feComposite> SVG filter primitive performs the combination of two input
// images pixel-wise in image space using one of the Porter-Duff compositing
// operations: over, in, atop, out, xor, and lighter. Additionally, a
// component-wise arithmetic operation (with the result clamped between [0..1])
// can be applied.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feComposite
func FeComposite(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feComposite", SVGNamespace, children...)
}

// The <feConvolveMatrix> SVG filter primitive applies a matrix convolution
// filter effect. A convolution combines pixels in the input image with
// neighboring pixels to produce a resulting image. A wide variety of imaging
// operations can be achieved through convolutions, including blurring, edge
// detection, sharpening, embossing and beveling.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feConvolveMatrix
func FeConvolveMatrix(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feConvolveMatrix", SVGNamespace, children...)
}

// The <feDiffuseLighting> SVG filter primitive lights an image using the alpha
// channel as a bump map. The resulting image, which is an RGBA opaque image,
// depends on the light color, light position and surface geometry of the input
// bump map.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feDiffuseLighting
func FeDiffuseLighting(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feDiffuseLighting", SVGNamespace, children...)
}

// The <feDisplacementMap> SVG filter primitive uses the pixel values from the
// image from in2 to spatially displace the image from in.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feDisplacementMap
func FeDisplacementMap(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feDisplacementMap", SVGNamespace, children...)
}

// The <feDistantLight> filter primitive defines a distant light source that
// can be used within a lighting filter primitive: <feDiffuseLighting> or
// <feSpecularLighting>.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feDistantLight
func FeDistantLight(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feDistantLight", SVGNamespace, children...)
}

// The <feFlood> SVG filter primitive fills the filter subregion with the color
// and opacity defined by flood-color and flood-opacity.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFlood
func FeFlood(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feFlood", SVGNamespace, children...)
}

// The <feFuncA> SVG filter primitive defines the transfer function for the
// alpha component of the input graphic of its parent <feComponentTransfer>
// element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFuncA
func FeFuncA(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feFuncA", SVGNamespace, children...)
}

// The <feFuncB> SVG filter primitive defines the transfer function for the
// blue component of the input graphic of its parent <feComponentTransfer>
// element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFuncB
func FeFuncB(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feFuncB", SVGNamespace, children...)
}

// The <feFuncG> SVG filter primitive defines the transfer function for the
// green component of the input graphic of its parent <feComponentTransfer>
// element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFuncG
func FeFuncG(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feFuncG", SVGNamespace, children...)
}

// The <feFuncR> SVG filter primitive defines the transfer function for the red
// component of the input graphic of its parent <feComponentTransfer> element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFuncR
func FeFuncR(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feFuncR", SVGNamespace, children...)
}

// The <feGaussianBlur> SVG filter primitive blurs the input image by the
// amount specified in stdDeviation, which defines the bell-curve.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feGaussianBlur
func FeGaussianBlur(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feGaussianBlur", SVGNamespace, children...)
}

// The <feImage> SVG filter primitive fetches image data from an external
// source and provides the pixel data as output (meaning if the external source
// is an SVG image, it is rasterized.)
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feImage
func FeImage(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feImage", SVGNamespace, children...)
}

// FeMerge allows filter effects to be applied concurrently instead of
// sequentially. This is achieved by other filters storing their output via the
// result attribute and then accessing it in a <feMergeNode> child.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feMerge
func FeMerge(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feMerge", SVGNamespace, children...)
}

// The feMergeNode takes the result of another filter to be processed by its
// parent <feMerge>.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feMergeNode
func FeMergeNode(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feMergeNode", SVGNamespace, children...)
}

// The <feMorphology> SVG filter primitive is used to erode or dilate the input
// image. It's usefulness lies especially in fattening or thinning effects.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feMorphology
func FeMorphology(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feMorphology", SVGNamespace, children...)
}

// The <feOffset> SVG filter primitive allows to offset the input image. The
// input image as a whole is offset by the values specified in the dx and dy
// attributes.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feOffset
func FeOffset(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feOffset", SVGNamespace, children...)
}

// The <fePointLight> filter primitive defines a light source which allows to
// create a point light effect. It that can be used within a lighting filter
// primitive: <feDiffuseLighting> or <feSpecularLighting>.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/fePointLight
func FePointLight(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("fePointLight", SVGNamespace, children...)
}

// The <feSpecularLighting> SVG filter primitive lights a source graphic using
// the alpha channel as a bump map. The resulting image is an RGBA image based
// on the light color. The lighting calculation follows the standard specular
// component of the Phong lighting model. The resulting image depends on the
// light color, light position and surface geometry of the input bump map. The
// result of the lighting calculation is added. The filter primitive assumes
// that the viewer is at infinity in the z direction.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feSpecularLighting
func FeSpecularLighting(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feSpecularLighting", SVGNamespace, children...)
}

// The <feSpotLight> SVG filter primitive defines a light source which allows
// to create a spotlight effect. It that can be used within a lighting filter
// primitive: <feDiffuseLighting> or <feSpecularLighting>.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feSpotLight
func FeSpotLight(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feSpotLight", SVGNamespace, children...)
}

// The <feTile> SVG filter primitive allows to fill a target rectangle with a
// repeated, tiled pattern of an input image. The effect is similar to the one
// of a <pattern>.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feTile
func FeTile(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feTile", SVGNamespace, children...)
}

// The <feTurbulence> SVG filter primitive creates an image using the Perlin
// turbulence function. It allows the synthesis of artificial textures like
// clouds or marble. The resulting image will fill the entire filter primitive
// subregion.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feTurbulence
func FeTurbulence(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("feTurbulence", SVGNamespace, children...)
}

// Filter serves as container for atomic filter operations. It is never
// rendered directly. A filter is referenced by using the filter attribute on
// the target SVG element or via the filter CSS property.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/filter
func Filter(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("filter", SVGNamespace, children...)
}

// Font defines a font to be used for text layout.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/font
func Font(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("font", SVGNamespace, children...)
}

// FontFace corresponds to the CSS @font-face rule. It defines a font's outer
// properties.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/font-face
func FontFace(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("font-face", SVGNamespace, children...)
}

// FontFaceFormat describes the type of font referenced by its parent
// <font-face-uri>.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/font-face-format
func FontFaceFormat(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("font-face-format", SVGNamespace, children...)
}

// The <font-face-name> element points to a locally installed copy of this
// font, identified by its name.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/font-face-name
func FontFaceName(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("font-face-name", SVGNamespace, children...)
}

// FontFaceSrc corresponds to the src descriptor in CSS @font-face rules. It
// serves as container for <font-face-name>, pointing to locally installed
// copies of this font, and <font-face-uri>, utilizing remotely defined fonts.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/font-face-src
func FontFaceSrc(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("font-face-src", SVGNamespace, children...)
}

// FontFaceUri points to a remote definition of the current font.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/font-face-uri
func FontFaceUri(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("font-face-uri", SVGNamespace, children...)
}

// ForeignObject allows for inclusion of a different XML namespace. In the
// context of a browser it is most likely XHTML/HTML.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/foreignObject
func ForeignObject(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("foreignObject", SVGNamespace, children...)
}

// G is a container used to group other SVG elements.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/g
func G(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("g", SVGNamespace, children...)
}

// A <glyph> defines a single glyph in an SVG font.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/glyph
func Glyph(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("glyph", SVGNamespace, children...)
}

// The glyphRef element provides a single possible glyph to the referencing
// <altGlyph> substitution.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/glyphRef
func GlyphRef(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("glyphRef", SVGNamespace, children...)
}

// Hkern allows to fine-tweak the horizontal distance between two glyphs. This
// process is known as kerning.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/hkern
func Hkern(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("hkern", SVGNamespace, children...)
}

// Image includes images inside SVG documents. It can display raster image
// files or other SVG files.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/image
func Image(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("image", SVGNamespace, children...)
}

// The <line> element is an SVG basic shape used to create a line connecting
// two points.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/line
func Line(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("line", SVGNamespace, children...)
}

// The <linearGradient> element lets authors define linear gradients that can
// be applied to fill or stroke of graphical elements.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/linearGradient
func LinearGradient(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("linearGradient", SVGNamespace, children...)
}

// The <marker> element defines the graphic that is to be used for drawing
// arrowheads or polymarkers on a given <path>, <line>, <polyline> or <polygon>
// element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/marker
func Marker(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("marker", SVGNamespace, children...)
}

// The <mask> element defines an alpha mask for compositing the current object
// into the background. A mask is used/referenced using the mask property.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/mask
func Mask(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("mask", SVGNamespace, children...)
}

// Metadata allows to add metadata to SVG content. Metadata is structured
// information about data. The contents of <metadata> elements should be
// elements from other XML namespaces such as RDF, FOAF, etc.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/metadata
func Metadata(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("metadata", SVGNamespace, children...)
}

// MissingGlyph's content is rendered, if for a given character the font
// doesn't define an appropriate <glyph>.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/missing-glyph
func MissingGlyph(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("missing-glyph", SVGNamespace, children...)
}

// The <mpath> sub-element for the <animateMotion> element provides the ability
// to reference an external <path> element as the definition of a motion path.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/mpath
func Mpath(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("mpath", SVGNamespace, children...)
}

// Path is the generic element to define a shape. All the basic shapes can be
// created with a path element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/path
func Path(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("path", SVGNamespace, children...)
}

// The <pattern> element defines a graphics object which can be redrawn at
// repeated x and y-coordinate intervals ("tiled") to cover an area.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/pattern
func Pattern(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("pattern", SVGNamespace, children...)
}

// The <polygon> element defines a closed shape consisting of a set of
// connected straight line segments. The last point is connected to the first
// point. For open shapes see the <polyline> element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/polygon
func Polygon(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("polygon", SVGNamespace, children...)
}

// Polyline is an SVG basic shape that creates straight lines connecting
// several points. Typically a polyline is used to create open shapes as the
// last point doesn't have to be connected to the first point. For closed
// shapes see the <polygon> element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/polyline
func Polyline(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("polyline", SVGNamespace, children...)
}

// RadialGradient lets authors define radial gradients to fill or stroke
// graphical elements.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/radialGradient
func RadialGradient(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("radialGradient", SVGNamespace, children...)
}

// The <rect> element is a basic SVG shape that creates rectangles, defined by
// their corner's position, their width, and their height. The rectangles may
// have their corners rounded.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/rect
func Rect(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("rect", SVGNamespace, children...)
}

// A SVG script element is equivalent to the script element in HTML and thus is
// the place for scripts (e.g., ECMAScript).
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/script
func Script(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("script", SVGNamespace, children...)
}

// The <set> element provides a simple means of just setting the value of an
// attribute for a specified duration. It supports all attribute types,
// including those that cannot reasonably be interpolated, such as string and
// boolean values. The <set> element is non-additive. The additive and
// accumulate attributes are not allowed, and will be ignored if specified.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/set
func Set(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("set", SVGNamespace, children...)
}

// Stop defines the ramp of colors to use on a gradient, which is a child
// element to either the <linearGradient> or the <radialGradient> element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/stop
func Stop(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("stop", SVGNamespace, children...)
}

// Style allows style sheets to be embedded directly within SVG content. SVG's
// style element has the same attributes as the corresponding element in HTML
// (see HTML's <style> element).
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/style
func Style(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("style", SVGNamespace, children...)
}

// The svg element is a container that defines a new coordinate system and
// viewport. It is used as the outermost element of any SVG document but it can
// also be used to embed a SVG fragment inside any SVG or HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/svg
func Svg(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("svg", SVGNamespace, children...)
}

// Switch evaluates the requiredFeatures, requiredExtensions and systemLanguage
// attributes on its direct child elements in order, and then processes and
// renders the first child for which these attributes evaluate to true. All
// others will be bypassed and therefore not rendered. If the child element is
// a container element such as a <g>, then the entire subtree is either
// processed/rendered or bypassed/not rendered.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/switch
func Switch(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("switch", SVGNamespace, children...)
}

// The <symbol> element is used to define graphical template objects which can
// be instantiated by a <use> element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/symbol
func Symbol(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("symbol", SVGNamespace, children...)
}

// The SVG <text> element defines a graphics element consisting of text. It's
// possible to apply a gradient, pattern, clipping path, mask, or filter to
// <text>, just like any other SVG graphics element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/text
func Text(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("text", SVGNamespace, children...)
}

// In addition to text drawn in a straight line, SVG also includes the ability
// to place text along the shape of a <path> element. To specify that a block
// of text is to be rendered along the shape of a <path>, include the given
// text within a <textPath> element which includes an href attribute with a
// reference to a <path> element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/textPath
func TextPath(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("textPath", SVGNamespace, children...)
}

// Each container element or graphics element in an SVG drawing can supply a
// <title> element containing a description string where the description is
// text-only. When the current SVG document fragment is rendered as SVG on
// visual media, <title> element is not rendered as part of the graphics.
// However, some user agents may, for example, display the <title> element as a
// tooltip. Alternate presentations are possible, both visual and aural, which
// display the <title> element but do not display path elements or other
// graphics elements. The <title> element generally improves accessibility of
// SVG documents.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/title
func Title(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("title", SVGNamespace, children...)
}

// The textual content for a <text> SVG element can be either character data
// directly embedded within the <text> element or the character data content of
// a referenced element, where the referencing is specified with a <tref>
// element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/tref
func Tref(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("tref", SVGNamespace, children...)
}

// Within a <text> element, text and font properties and the current text
// position can be adjusted with absolute or relative coordinate values by
// including a <tspan> element.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/tspan
func Tspan(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("tspan", SVGNamespace, children...)
}

// The <use> element takes nodes from within the SVG document, and duplicates
// them somewhere else.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/use
func Use(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("use", SVGNamespace, children...)
}

// A view is a defined way to view the image, like a zoom level or a detail
// view.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/view
func View(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("view", SVGNamespace, children...)
}

// Vkern allows to fine-tweak the vertical distance between two glyphs in
// top-to-bottom fonts. This process is known as kerning.
//
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/vkern
func Vkern(children ...markup.MarkupOrChild) markup.RenderResult {
	return markup.TagWithNS("vkern", SVGNamespace, children...)
}
