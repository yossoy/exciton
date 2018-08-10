package markup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"golang.org/x/net/html/atom"

	"golang.org/x/net/html"
)

func resolvePathNode(root *html.Node, path []int) *html.Node {
	c := root
	for _, p := range path {
		c = c.FirstChild
		for p > 0 {
			c = c.NextSibling
			p = p - 1
		}
	}
	return c
}

func getAttrFromNode(n *html.Node, name string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == name {
			return attr.Val, true
		}
	}
	return "", false
}

func setAttrToNode(n *html.Node, name string, val string) {
	bProc := false
	for idx, attr := range n.Attr {
		if attr.Key == name {
			n.Attr[idx].Val = val
			bProc = true
			break
		}
	}
	if !bProc {
		n.Attr = append(n.Attr, html.Attribute{
			Key: name,
			Val: val,
		})
	}
}

func delAttrFromNode(n *html.Node, name string) {
	for idx, attr := range n.Attr {
		if attr.Key == name {
			n.Attr = append(n.Attr[:idx], n.Attr[idx+1:]...)
			break
		}
	}
}

func makeEventName(id string, pd bool, sp bool) string {
	return fmt.Sprintf("evt{%s,%v,%v}", id, pd, sp)
}

func applyDiff(ds *DiffSet, root *html.Node) *html.Node {
	var (
		curNode  *html.Node
		arg1Node *html.Node
		arg2Node *html.Node
		creNodes []*html.Node
	)
	retNode := root
	for _, itm := range ds.Items {
		switch itm.ItemType {
		case ditCreateNode:
			n := itm.Value.(string)
			curNode = &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Lookup([]byte(n)),
				Data:     n,
			}
			creNodes = append(creNodes, curNode)
		case ditCreateNodeWithNS:
			n := itm.Value.(string)
			curNode = &html.Node{
				Type:      html.ElementNode,
				DataAtom:  atom.Lookup([]byte(n)),
				Data:      n,
				Namespace: itm.Key,
			}
			creNodes = append(creNodes, curNode)
		case ditCreateTextNode:
			n := itm.Value.(string)
			curNode = &html.Node{
				Type: html.TextNode,
				Data: n,
			}
			creNodes = append(creNodes, curNode)
		case ditSelectCurNode:
			if idx, ok := itm.Value.(int); ok {
				curNode = creNodes[idx]
			} else if p, ok := itm.Value.([]int); ok {
				curNode = resolvePathNode(root, p) //p[1:])
			} else {
				panic(itm.Value)
			}
		case ditSelectArg1Node:
			if idx, ok := itm.Value.(int); ok {
				arg1Node = creNodes[idx]
			} else if p, ok := itm.Value.([]int); ok {
				arg1Node = resolvePathNode(root, p) //p[1:])
			} else {
				panic(itm.Value)
			}
		case ditSelectArg2Node:
			if idx, ok := itm.Value.(int); ok {
				arg2Node = creNodes[idx]
			} else if p, ok := itm.Value.([]int); ok {
				arg2Node = resolvePathNode(root, p) //p[1:])
			} else {
				panic(itm.Value)
			}
		case ditPropertyValue:
			k := itm.Key
			v := itm.Value
			setAttrToNode(curNode, "_prop_"+k, fmt.Sprintf("%v", v))
		case ditDelProperty:
			n := itm.Value.(string)
			delAttrFromNode(curNode, "_prop_"+n)
		case ditAttributeValue:
			k := itm.Key
			v := itm.Value
			setAttrToNode(curNode, k, v.(string))
		case ditDelAttributeValue:
			n := itm.Value.(string)
			delAttrFromNode(curNode, n)
		case ditAddClassList:
			k := itm.Value.(string)
			cl, ok := getAttrFromNode(curNode, "class")
			if ok {
				for _, s := range strings.Split(cl, " ") {
					if s == k {
						panic(fmt.Errorf("already exists in class: classes:%q, class:%q", cl, k))
					}
				}
				cl = cl + " "
			}
			setAttrToNode(curNode, "class", cl+k)
		case ditDelClassList:
			k := itm.Value.(string)
			if cl, ok := getAttrFromNode(curNode, "class"); ok {
				cll := strings.Split(cl, " ")
				found := false
				for idx, s := range cll {
					if s == k {
						found = true
						cll = append(cll[:idx], cll[idx+1:]...)
						break
					}
				}
				if !found {
					panic(fmt.Errorf("not exists in class: classes:%q, class:%q", cl, k))
				}
				if len(cll) == 0 {
					delAttrFromNode(curNode, "class")
				} else {
					setAttrToNode(curNode, "class", strings.Join(cll, " "))
				}
			}
		case ditAddDataSet:
			k := itm.Key
			v := itm.Value.(string)
			setAttrToNode(curNode, datasetToAttr(k), v)
		case ditDelDataSet:
			n := itm.Value.(string)
			delAttrFromNode(curNode, datasetToAttr(n))
		case ditAddStyle:
			k := itm.Key
			v := itm.Value.(string)
			style := make(map[string]string)
			if s, ok := getAttrFromNode(curNode, "style"); ok {
				json.Unmarshal([]byte(s), &style)
			}
			style[k] = v
			if s, err := json.Marshal(&style); err == nil {
				setAttrToNode(curNode, "style", string(s))
			}
		case ditDelStyle:
			n := itm.Value.(string)
			if s, ok := getAttrFromNode(curNode, "style"); ok {
				style := make(map[string]string)
				json.Unmarshal([]byte(s), &style)
				delete(style, n)
				if len(style) == 0 {
					delAttrFromNode(curNode, "style")
				} else {
					if b, err := json.Marshal(&style); err == nil {
						setAttrToNode(curNode, "style", string(b))
					}
				}
			}
		case ditNodeValue:
			if curNode.Type != html.TextNode {
				panic("invalid target")
			}
			s := itm.Value.(string)
			curNode.Data = s
		case ditInnerHTML:
			v := itm.Value.(string)
			nn, err := html.ParseFragment(strings.NewReader(v), curNode)
			if err != nil {
				panic(err)
			}
			for {
				c := curNode.FirstChild
				if c == nil {
					break
				}
				curNode.RemoveChild(c)
			}
			for _, n := range nn {
				curNode.AppendChild(n)
			}
		case ditAppendChild:
			if arg1Node.Parent != nil {
				arg1Node.Parent.RemoveChild(arg1Node)
			}
			curNode.AppendChild(arg1Node)
		case ditInsertBefore:
			curNode.InsertBefore(arg1Node, arg2Node)
		case ditRemoveChild:
			curNode.RemoveChild(arg1Node)
		case ditReplaceChild:
			//curNode.ReplaceChild(arg1Node, arg2Node)
			if arg2Node.Parent != curNode {
				panic("invalid parent")
			}
			if arg1Node.Parent != nil {
				arg1Node.Parent.RemoveChild(arg1Node)
			}
			curNode.InsertBefore(arg1Node, arg2Node)
			curNode.RemoveChild(arg2Node)
		case ditAddEventListener:
			k := itm.Key
			v := itm.Value.(dtAddEventListenerItem)
			setAttrToNode(curNode, "on"+k, makeEventName(v.ID, v.PreventDefault, v.StopPropagation))
		case ditRemoveEventListener:
			k := itm.Value.(string)
			delAttrFromNode(curNode, "on"+k)
		case ditSetRootItem:
			retNode = curNode
		case ditNodeUUID:
			v := itm.Value
			setAttrToNode(curNode, "_uuid", v.(string))
		case ditMountComponent:
			//TODO:
		case ditUnmountComponent:
			//TODO:
		default:
			panic("invalid input")
		}
	}

	return retNode
}

func loadHTML(ds *DiffSet, h string) *node {
	var curNode *node
	var rootNode *node
	z := html.NewTokenizer(strings.NewReader(h))
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() != io.EOF {
				panic(z.Err())
			}
			return rootNode
		case html.SelfClosingTagToken:
			fallthrough
		case html.StartTagToken:
			n, hasAttr := z.TagName()
			tagName := string(n)
			nn := &node{tag: tagName}
			ds.createNode(nn)
			if hasAttr {
				for {
					k, v, m := z.TagAttr()
					ds.AddAttribute(nn, string(k), string(v))
					if !m {
						break
					}
				}
			}
			if curNode != nil {
				ds.appendChild(curNode, nn)
				curNode.appendChild(nn)
			}
			if tt == html.StartTagToken {
				curNode = nn
			}
			if rootNode == nil {
				rootNode = nn
			}
		case html.EndTagToken:
			if curNode == nil {
				panic("parse fail")
			}
			curNode = curNode.parent
		case html.TextToken:
			t := string(z.Text())
			nn := &node{text: t}
			ds.createTextNode(nn)
			if curNode != nil {
				ds.appendChild(curNode, nn)
				curNode.appendChild(nn)
			}
		}
	}
}

func TestDiffHTML(t *testing.T) {
	rootNode := &node{rootNode: true}
	ds := DiffSet{rootNode: rootNode}
	h := `<div id="id1">aaa</div>`
	rn := loadHTML(&ds, h)
	//ds.SetRootItem(rn, "root")
	ds.appendChild(rootNode, rn)
	rn.parent = rootNode
	//	rn.setRootNode()
	b, _ := json.Marshal(ds)
	t.Logf("%s:(%d bytes)", string(b), len(b))

	//var root *html.Node
	root := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Lookup([]byte("div")),
		Data:     "div",
	}
	//root =
	applyDiff(&ds, root)
	buf := bytes.NewBufferString("")
	html.Render(buf, root)
	t.Log(buf.String())
	if buf.String() != "<div>"+h+"</div>" {
		t.Errorf("invalid apply[1] result: %q vs %q", h, buf.String())
	}

	ds.reset()

	// insert <a href="http://foo.bar>anchor"</a> before aaa
	n := &node{tag: "a"}
	ds.createNode(n)
	ds.AddAttribute(n, "href", "http://foo.bar")
	nt := &node{text: "anchor"}
	ds.createTextNode(nt)
	ds.appendChild(n, nt)
	n.appendChild(nt)
	ds.insertBefore(rn, n, rn.firstChild())
	//ds.appendChild(rn, n)
	rn.insertBefore(n, rn.firstChild())
	//rn.appendChild(n)
	b, _ = json.Marshal(ds)
	t.Logf("%s:(%d bytes)", string(b), len(b))
	root = applyDiff(&ds, root)
	buf = bytes.NewBufferString("")
	html.Render(buf, root)
	t.Log(buf.String())
	h = `<div><div id="id1"><a href="http://foo.bar">anchor</a>aaa</div></div>`
	if h != buf.String() {
		t.Errorf("invalid apply[2] result: %q vs %q", h, buf.String())
	}

}

func TestDiffItem1(t *testing.T) {
	rootNode := &node{rootNode: true}
	ds := DiffSet{rootNode: rootNode}
	n1 := &node{tag: "A"}
	ds.createNode(n1)
	ds.AddAttribute(n1, "a", "b")
	n2 := &node{tag: "B"}
	ds.createNode(n2)
	ds.addProperty(n2, "b", "c")
	ds.appendChild(n1, n2)
	n1.appendChild(n2)
	ds.SetRootItem(n1, "root")
	//n1.setRootNode()
	n1.parent = rootNode
	b, _ := json.MarshalIndent(&ds, "  ", "  ")
	t.Log(string(b))
	ds.reset()
	ds.AddAttribute(n2, "a", "c")
	n3 := &node{ns: "aa", tag: "bb"}
	ds.createNodeWithNS(n3)
	ds.appendChild(n2, n3)
	n2.appendChild(n3)
	b, _ = json.MarshalIndent(&ds, "", "  ")
	t.Log(string(b))
}
