package markup

import (
	"reflect"
	"strings"
)

type delayApplyer func(b Builder) interface{}

type attrApplyer struct {
	name      string
	nameSpace string
	value     interface{}
}

func (aa attrApplyer) isMarkup()        {}
func (aa attrApplyer) isMarkupOrChild() {}
func (aa attrApplyer) applyToNode(b Builder, n Node, on Node) {
	bb := b.(*builder)
	nn := n.(*node)
	onn := on.(*node)
	k := aa.nameSpace + ":" + aa.name
	ov, ok := onn.attributes[k]
	if ok {
		delete(onn.attributes, k)
	}
	if nn.attributes == nil {
		nn.attributes = make(map[string]interface{})
	}
	val := aa.value
	if da, ok := val.(delayApplyer); ok {
		val = da(b)
	}
	nn.attributes[k] = val
	if !ok || ov != val {
		if nn.ns == "" && aa.nameSpace == "" {
			bb.diffSet.AddAttribute(nn, aa.name, val)
		} else {
			var ns interface{}
			if aa.nameSpace != "" {
				ns = aa.nameSpace
			}
			bb.diffSet.AddAttributeNS(nn, aa.name, ns, val)
		}
	}
}

func attrToDataSet(attr string) (string, bool) {
	if !strings.HasPrefix(attr, "data-") {
		return "", false
	}
	rr := make([]rune, 0, len(attr))
	prevdash := false
	for i, r := range attr {
		if i < len("data-") {
			// skip "data-" prefix
			continue
		}
		//TODO: validate rune range?
		if r == '-' {
			prevdash = true
			continue
		}
		if prevdash && ('a' <= r) && (r <= 'z') {
			r = r - 'a' + 'A'
		}
		rr = append(rr, r)
		prevdash = false
	}
	return string(rr), true
}

func datasetToAttr(ds string) string {
	rr := make([]rune, 0, len(ds))
	for i, r := range ds {
		if (i != 0) && ('A' <= r) && (r <= 'Z') {
			rr = append(rr, '-')
			r = r - 'A' + 'a'
		}
		rr = append(rr, r)
	}
	return "data-" + string(rr)
}

func Attribute(name string, value interface{}) MarkupOrChild {
	//if name is start "data-", change name and return dataApplyer.
	if ds, ok := attrToDataSet(name); ok {
		return Data(ds, value.(string))
	}

	return attrApplyer{
		name:  name,
		value: value,
	}
}

func AttributeNS(name string, nameSpace string, value interface{}) MarkupOrChild {
	return attrApplyer{
		name:      name,
		nameSpace: nameSpace,
		value:     value,
	}
}

type propApplyer struct {
	name       string
	value      interface{}
	isRedirect bool
}

func (aa propApplyer) isMarkup()        {}
func (aa propApplyer) isMarkupOrChild() {}
func (aa propApplyer) applyToNode(b Builder, n Node, on Node) {
	bb := b.(*builder)
	nn := n.(*node)
	onn := on.(*node)
	ov, ok := onn.properties[aa.name]
	if ok {
		delete(onn.properties, aa.name)
	}
	if nn.properties == nil {
		nn.properties = make(map[string]interface{})
	}
	val := aa.value
	if da, ok := val.(delayApplyer); ok {
		val = da(b)
	}
	nn.properties[aa.name] = val
	if !ok || ov != val {
		bb.diffSet.addProperty(nn, aa.name, val)
	}
}
func (aa propApplyer) applyToComponent(c Component) {
	core := c.Context()
	if idx, ok := core.klass.Properties[aa.name]; ok {
		v := reflect.ValueOf(c)
		val := aa.value
		if da, ok := val.(delayApplyer); ok {
			val = da(c.Builder())
		}
		vv := reflect.ValueOf(val)
		v.Elem().Field(idx).Set(vv)
	}
}

func Property(name string, value interface{}) MarkupOrChild {
	return propApplyer{
		name:  name,
		value: value,
	}
}

type dataApplyer struct {
	name  string
	value string
}

func (da dataApplyer) isMarkup()        {}
func (da dataApplyer) isMarkupOrChild() {}
func (da dataApplyer) applyToNode(b Builder, n Node, on Node) {
	bb := b.(*builder)
	nn := n.(*node)
	onn := on.(*node)
	ov, ok := onn.dataset[da.name]
	if ok {
		delete(onn.dataset, da.name)
	}
	if nn.dataset == nil {
		nn.dataset = make(map[string]string)
	}
	nn.dataset[da.name] = da.value
	if !ok || ov != da.value {
		bb.diffSet.AddDataSet(nn, da.name, da.value)
	}
}

func Data(name string, value string) MarkupOrChild {
	return dataApplyer{
		name:  name,
		value: value,
	}
}

type classApplyer []string

func (ca classApplyer) isMarkup()        {}
func (ca classApplyer) isMarkupOrChild() {}
func (ca classApplyer) applyToNode(b Builder, n Node, on Node) {
	bb := b.(*builder)
	nn := n.(*node)
	onn := on.(*node)
	for _, c := range ca {
		_, ok := onn.classes[c]
		if ok {
			delete(onn.classes, c)
		}
		if nn.classes == nil {
			nn.classes = make(map[string]struct{})
		}
		nn.classes[c] = struct{}{}
		if !ok {
			bb.diffSet.AddClassList(nn, c)
		}
	}
}

func Classes(classes ...string) MarkupOrChild {
	return classApplyer(classes)
}

type styleApplyer struct {
	name  string
	value string
}

func (sa styleApplyer) isMarkup()        {}
func (sa styleApplyer) isMarkupOrChild() {}
func (sa styleApplyer) applyToNode(b Builder, n Node, on Node) {
	bb := b.(*builder)
	nn := n.(*node)
	onn := on.(*node)
	ov, ok := onn.styles[sa.name]
	if ok {
		delete(onn.styles, sa.name)
	}
	if nn.styles == nil {
		nn.styles = make(map[string]string)
	}
	nn.styles[sa.name] = sa.value
	if !ok || ov != sa.value {
		bb.diffSet.AddStyle(nn, sa.name, sa.value)
	}
}

func Style(name string, value string) MarkupOrChild {
	return styleApplyer{
		name:  name,
		value: value,
	}
}

type innerHTMLApplyer string

func (iha innerHTMLApplyer) isMarkup()        {}
func (iha innerHTMLApplyer) isMarkupOrChild() {}
func (iha innerHTMLApplyer) applyToNode(b Builder, n Node, on Node) {
	nn := n.(*node)
	onn := on.(*node)
	nv := string(iha)
	ov := onn.innerHTML
	nn.innerHTML = nv
	if ov != nv {
		b.(*builder).diffSet.AddInnerHTML(nn, nv)
	}
}
func UnsafeHTML(html string) MarkupOrChild {
	return innerHTMLApplyer(html)
}
