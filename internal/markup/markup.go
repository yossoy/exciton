package markup

import (
	"reflect"
	"strings"
)

type delayApplyer func(b Builder) interface{}

type AttrApplyer struct {
	Name      string
	NameSpace string
	Value     interface{}
}

func (aa AttrApplyer) isMarkup()        {}
func (aa AttrApplyer) isMarkupOrChild() {}
func (aa AttrApplyer) applyToNode(b Builder, n Node, on Node) {
	bb := b.(*builder)
	nn := n.(*node)
	onn := on.(*node)
	k := aa.NameSpace + ":" + aa.Name
	ov, ok := onn.attributes[k]
	if ok {
		delete(onn.attributes, k)
	}
	if nn.attributes == nil {
		nn.attributes = make(map[string]interface{})
	}
	val := aa.Value
	if da, ok := val.(delayApplyer); ok {
		val = da(b)
	}
	nn.attributes[k] = val
	if !ok || ov != val {
		if nn.ns == "" && aa.NameSpace == "" {
			bb.diffSet.AddAttribute(nn, aa.Name, val)
		} else {
			var ns interface{}
			if aa.NameSpace != "" {
				ns = aa.NameSpace
			}
			bb.diffSet.AddAttributeNS(nn, aa.Name, ns, val)
		}
	}
}

func AttrToDataSet(attr string) (string, bool) {
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

func DatasetToAttr(ds string) string {
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

type PropApplyer struct {
	Name       string
	Value      interface{}
	IsRedirect bool
}

func (aa PropApplyer) isMarkup()        {}
func (aa PropApplyer) isMarkupOrChild() {}
func (aa PropApplyer) applyToNode(b Builder, n Node, on Node) {
	bb := b.(*builder)
	nn := n.(*node)
	onn := on.(*node)
	ov, ok := onn.properties[aa.Name]
	if ok {
		delete(onn.properties, aa.Name)
	}
	if nn.properties == nil {
		nn.properties = make(map[string]interface{})
	}
	val := aa.Value
	if da, ok := val.(delayApplyer); ok {
		val = da(b)
	}
	nn.properties[aa.Name] = val
	if !ok || ov != val {
		bb.diffSet.addProperty(nn, aa.Name, val)
	}
}
func (aa PropApplyer) applyToComponent(c Component) {
	core := c.Context()
	if idx, ok := core.klass.Properties[aa.Name]; ok {
		v := reflect.ValueOf(c)
		val := aa.Value
		if da, ok := val.(delayApplyer); ok {
			val = da(c.Builder())
		}
		vv := reflect.ValueOf(val)
		v.Elem().Field(idx).Set(vv)
	}
}

type DataApplyer struct {
	Name  string
	Value string
}

func (da DataApplyer) isMarkup()        {}
func (da DataApplyer) isMarkupOrChild() {}
func (da DataApplyer) applyToNode(b Builder, n Node, on Node) {
	bb := b.(*builder)
	nn := n.(*node)
	onn := on.(*node)
	ov, ok := onn.dataset[da.Name]
	if ok {
		delete(onn.dataset, da.Name)
	}
	if nn.dataset == nil {
		nn.dataset = make(map[string]string)
	}
	nn.dataset[da.Name] = da.Value
	if !ok || ov != da.Value {
		bb.diffSet.AddDataSet(nn, da.Name, da.Value)
	}
}

type ClassApplyer []string

func (ca ClassApplyer) isMarkup()        {}
func (ca ClassApplyer) isMarkupOrChild() {}
func (ca ClassApplyer) applyToNode(b Builder, n Node, on Node) {
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

type StyleApplyer struct {
	Name  string
	Value string
}

func (sa StyleApplyer) isMarkup()        {}
func (sa StyleApplyer) isMarkupOrChild() {}
func (sa StyleApplyer) applyToNode(b Builder, n Node, on Node) {
	bb := b.(*builder)
	nn := n.(*node)
	onn := on.(*node)
	ov, ok := onn.styles[sa.Name]
	if ok {
		delete(onn.styles, sa.Name)
	}
	if nn.styles == nil {
		nn.styles = make(map[string]string)
	}
	nn.styles[sa.Name] = sa.Value
	if !ok || ov != sa.Value {
		bb.diffSet.AddStyle(nn, sa.Name, sa.Value)
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
