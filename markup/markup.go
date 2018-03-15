package markup

import (
	"reflect"
	"strings"
)

type attrApplyer struct {
	name  string
	value interface{}
}

func (aa attrApplyer) isMarkup()        {}
func (aa attrApplyer) isMarkupOrChild() {}
func (aa attrApplyer) applyToNode(b *Builder, n *node, on *node) {
	ov, ok := on.attributes[aa.name]
	if ok {
		delete(on.attributes, aa.name)
	}
	if n.attributes == nil {
		n.attributes = make(map[string]interface{})
	}
	n.attributes[aa.name] = aa.value
	if !ok || ov != aa.value {
		b.diffSet.AddAttribute(n, aa.name, aa.value)
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

type propApplyer struct {
	name  string
	value interface{}
}

func (aa propApplyer) isMarkup()        {}
func (aa propApplyer) isMarkupOrChild() {}
func (aa propApplyer) applyToNode(b *Builder, n *node, on *node) {
	ov, ok := on.properties[aa.name]
	if ok {
		delete(on.properties, aa.name)
	}
	if n.properties == nil {
		n.properties = make(map[string]interface{})
	}
	n.properties[aa.name] = aa.value
	if !ok || ov != aa.value {
		b.diffSet.addProperty(n, aa.name, aa.value)
	}
}
func (aa propApplyer) applyToComponent(c Component) {
	core := c.Context()
	if idx, ok := core.klass.Properties[aa.name]; ok {
		v := reflect.ValueOf(c)
		vv := reflect.ValueOf(aa.value)
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
func (da dataApplyer) applyToNode(b *Builder, n *node, on *node) {
	ov, ok := on.dataset[da.name]
	if ok {
		delete(on.dataset, da.name)
	}
	if n.dataset == nil {
		n.dataset = make(map[string]string)
	}
	n.dataset[da.name] = da.value
	if !ok || ov != da.value {
		b.diffSet.AddDataSet(n, da.name, da.value)
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
func (ca classApplyer) applyToNode(b *Builder, n *node, on *node) {
	for _, c := range ca {
		_, ok := on.classes[c]
		if ok {
			delete(on.classes, c)
		}
		if n.classes == nil {
			n.classes = make(map[string]struct{})
		}
		n.classes[c] = struct{}{}
		if !ok {
			b.diffSet.AddClassList(n, c)
		}
	}
}

func Classes(class ...string) MarkupOrChild {
	return classApplyer(class)
}

type styleApplyer struct {
	name  string
	value string
}

func (sa styleApplyer) isMarkup()        {}
func (sa styleApplyer) isMarkupOrChild() {}
func (sa styleApplyer) applyToNode(b *Builder, n *node, on *node) {
	ov, ok := on.styles[sa.name]
	if ok {
		delete(on.styles, sa.name)
	}
	if n.styles == nil {
		n.styles = make(map[string]string)
	}
	n.styles[sa.name] = sa.value
	if !ok || ov != sa.value {
		b.diffSet.AddStyle(n, sa.name, sa.value)
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
func (iha innerHTMLApplyer) applyToNode(b *Builder, n *node, on *node) {
	nv := string(iha)
	ov := on.innerHTML
	n.innerHTML = nv
	if ov != nv {
		b.diffSet.AddInnerHTML(n, nv)
	}
}
func UnsafeHTML(html string) MarkupOrChild {
	return innerHTMLApplyer(html)
}
