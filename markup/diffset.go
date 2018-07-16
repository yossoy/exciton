package markup

import "github.com/yossoy/exciton/log"

type diffItemType int

const (
	ditNone diffItemType = iota
	ditCreateNode
	ditCreateNodeWithNS
	ditCreateTextNode
	ditSelectCurNode
	ditSelectArg1Node
	ditSelectArg2Node
	ditPropertyValue
	ditDelProperty
	ditAttributeValue
	ditDelAttributeValue
	ditAddClassList
	ditDelClassList
	ditAddDataSet
	ditDelDataSet
	ditAddStyle
	ditDelStyle
	ditNodeValue
	ditInnerHTML
	ditAppendChild
	ditInsertBefore
	ditRemoveChild
	ditReplaceChild
	ditAddEventListener
	ditRemoveEventListener
	ditSetRootItem
	ditNodeUUID
	ditAddClientEvent
	ditMountComponent
	ditUnmountComponent
)

func (t diffItemType) String() string {
	switch t {
	case ditNone:
		return "none"
	case ditCreateNode:
		return "createNode"
	case ditCreateNodeWithNS:
		return "createNodeNS"
	case ditCreateTextNode:
		return "createTextNode"
	case ditSelectCurNode:
		return "selectCurNode"
	case ditSelectArg1Node:
		return "selectArg1Node"
	case ditSelectArg2Node:
		return "selectArg2Node"
	case ditPropertyValue:
		return "propValue"
	case ditDelProperty:
		return "delProp"
	case ditAttributeValue:
		return "attrValue"
	case ditDelAttributeValue:
		return "delAttr"
	case ditAddClassList:
		return "addClass"
	case ditDelClassList:
		return "delClass"
	case ditAddDataSet:
		return "addDataSet"
	case ditDelDataSet:
		return "delDataSet"
	case ditAddStyle:
		return "addStyle"
	case ditDelStyle:
		return "delStyle"
	case ditNodeValue:
		return "nodeValue"
	case ditInnerHTML:
		return "innerHTML"
	case ditAppendChild:
		return "appendChild"
	case ditInsertBefore:
		return "insertBefore"
	case ditRemoveChild:
		return "removeChild"
	case ditReplaceChild:
		return "replaceChild"
	case ditAddEventListener:
		return "addEventListener"
	case ditRemoveEventListener:
		return "removeEventListener"
	case ditSetRootItem:
		return "setRootItem"
	case ditNodeUUID:
		return "nodeUUID"
	case ditAddClientEvent:
		return "addClientEvent"
	case ditMountComponent:
		return "mountComponent"
	case ditUnmountComponent:
		return "unmountComponent"
	default:
		panic("invalid diffType")
	}
}

type diffItem struct {
	ItemType diffItemType `json:"t"`
	//DebugString string       `json:"type,omitempty"`
	DebugString string      `json:"-"`
	Key         string      `json:"k,omitempty"`
	Value       interface{} `json:"v,omitempty"`
}

// DiffSet is store of node diff data
type DiffSet struct {
	Items    []diffItem `json:"items"`
	newNodes map[*node]int
	curNode  *node
	arg1Node *node
	arg2Node *node
	rootNode *node
}

func (ds *DiffSet) hasDiff() bool {
	return len(ds.Items) > 0
}

func (ds *DiffSet) reset() {
	ds.Items = ds.Items[:0]
	for k := range ds.newNodes {
		delete(ds.newNodes, k)
	}
	ds.curNode = nil
	ds.arg1Node = nil
	ds.arg2Node = nil
}

func (ds *DiffSet) addItem(t diffItemType, arg interface{}) {
	ds.Items = append(ds.Items, diffItem{
		ItemType:    t,
		DebugString: t.String(),
		Value:       arg,
	})
}

func (ds *DiffSet) addItemWithKey(t diffItemType, key string, arg interface{}) {
	ds.Items = append(ds.Items, diffItem{
		ItemType:    t,
		DebugString: t.String(),
		Key:         key,
		Value:       arg,
	})
}

func (ds *DiffSet) addNewNode(t diffItemType, arg *node) {
	idx := len(ds.newNodes)
	if ds.newNodes == nil {
		ds.newNodes = make(map[*node]int)
	}
	ds.newNodes[arg] = idx
	ds.curNode = arg
	switch t {
	case ditCreateNode:
		ds.addItem(t, arg.tag)
	case ditCreateNodeWithNS:
		ds.addItemWithKey(t, arg.ns, arg.tag)
	case ditCreateTextNode:
		ds.addItem(t, arg.text)
	default:
		panic("exciton: invalid argument")
	}
}

func (ds *DiffSet) selectCurNode(t *node) {
	switch {
	case ds.rootNode == t:
		ds.addItem(ditSelectCurNode, ds.rootNode.indexPath(ds.rootNode))
	case ds.curNode == t:
		return
	case t.isMount(ds.rootNode):
		ds.addItem(ditSelectCurNode, t.indexPath(ds.rootNode))
	default:
		if idx, ok := ds.newNodes[t]; ok {
			ds.addItem(ditSelectCurNode, idx)
		} else {
			log.PrintError("invalid node!: %#v\n", t.parent)
			panic("exciton: invalid sequence[unmounted node select]")
		}
	}
	ds.curNode = t
}

func (ds *DiffSet) selectArg1Node(t *node) {
	switch {
	case ds.rootNode == t:
		ds.addItem(ditSelectCurNode, ds.rootNode.indexPath(ds.rootNode))
	case ds.arg1Node == t:
		return
	case t.isMount(ds.rootNode):
		ds.addItem(ditSelectArg1Node, t.indexPath(ds.rootNode))
	default:
		if idx, ok := ds.newNodes[t]; ok {
			ds.addItem(ditSelectArg1Node, idx)
		} else {
			panic("exciton: invalid sequence[un mounted node select]")
		}
	}
	ds.arg1Node = t
}

func (ds *DiffSet) selectArg2Node(t *node) {
	switch {
	case ds.rootNode == t:
		ds.addItem(ditSelectCurNode, ds.rootNode.indexPath(ds.rootNode))
	case ds.arg2Node == t:
		return
	case t.isMount(ds.rootNode):
		ds.addItem(ditSelectArg2Node, t.indexPath(ds.rootNode))
	default:
		if idx, ok := ds.newNodes[t]; ok {
			ds.addItem(ditSelectArg2Node, idx)
		} else {
			panic("exciton: invalid sequence[un mounted node select]")
		}
	}
	ds.arg2Node = t
}

func (ds *DiffSet) createNode(t *node) {
	ds.addNewNode(ditCreateNode, t)
}

func (ds *DiffSet) createNodeWithNS(t *node) {
	ds.addNewNode(ditCreateNodeWithNS, t)
}

func (ds *DiffSet) createTextNode(t *node) {
	ds.addNewNode(ditCreateTextNode, t)
}

func (ds *DiffSet) setNodeUUID(t *node, id string) {
	ds.selectCurNode(t)
	ds.addItem(ditNodeUUID, id)
}

func (ds *DiffSet) setNodeValue(t *node, nodeValue string) {
	ds.selectCurNode(t)
	ds.addItem(ditNodeValue, nodeValue)
}

//TODO: merge
func (ds *DiffSet) addProperty(t *node, name string, value interface{}) {
	ds.selectCurNode(t)
	ds.addItemWithKey(ditPropertyValue, name, value)
}
func (ds *DiffSet) delProperty(t *node, name string) {
	ds.selectCurNode(t)
	ds.addItem(ditDelProperty, name)
}

func (ds *DiffSet) AddAttribute(t *node, name string, value interface{}) {
	ds.selectCurNode(t)
	ds.addItemWithKey(ditAttributeValue, name, value)
}
func (ds *DiffSet) DelAttribute(t *node, name string) {
	ds.selectCurNode(t)
	ds.addItem(ditDelAttributeValue, name)
}

func (ds *DiffSet) AddClassList(t *node, value string) {
	ds.selectCurNode(t)
	ds.addItem(ditAddClassList, value)
}
func (ds *DiffSet) DelClassList(t *node, value string) {
	ds.selectCurNode(t)
	ds.addItem(ditDelClassList, value)
}

func (ds *DiffSet) AddInnerHTML(t *node, value string) {
	ds.selectCurNode(t)
	ds.addItem(ditInnerHTML, value)
}

func (ds *DiffSet) AddDataSet(t *node, name string, value string) {
	ds.selectCurNode(t)
	ds.addItemWithKey(ditAddDataSet, name, value)
}
func (ds *DiffSet) DelDataSet(t *node, name string) {
	ds.selectCurNode(t)
	ds.addItem(ditDelDataSet, name)
}

func (ds *DiffSet) AddStyle(t *node, name string, value string) {
	ds.selectCurNode(t)
	ds.addItemWithKey(ditAddStyle, name, value)
}
func (ds *DiffSet) DelStyle(t *node, name string) {
	ds.selectCurNode(t)
	ds.addItem(ditDelStyle, name)
}

func (ds *DiffSet) appendChild(t *node, c *node) {
	ds.selectCurNode(t)
	ds.selectArg1Node(c)
	ds.addItem(ditAppendChild, nil)
}

func (ds *DiffSet) insertBefore(t *node, n *node, c *node) {
	ds.selectCurNode(t)
	ds.selectArg1Node(n)
	ds.selectArg2Node(c)
	ds.addItem(ditInsertBefore, nil)
}

func (ds *DiffSet) RemoveChild(t *node, c *node) {
	ds.selectCurNode(t)
	ds.selectArg1Node(c)
	ds.addItem(ditRemoveChild, nil)
}

func (ds *DiffSet) ReplaceChild(p *node, n *node, o *node) {
	ds.selectCurNode(p)
	ds.selectArg1Node(n)
	ds.selectArg2Node(o)
	ds.addItem(ditReplaceChild, nil)
}

type dtAddEventListenerItem struct {
	ID              string `json:"id"`
	PreventDefault  bool   `json:"pd"`
	StopPropagation bool   `json:"sp"`
}

func (ds *DiffSet) AddEventListener(t *node, name string, eventId string, preventDefault bool, stopPropagation bool) {
	ds.selectCurNode(t)
	ds.addItemWithKey(ditAddEventListener, name, dtAddEventListenerItem{
		ID:              eventId,
		PreventDefault:  preventDefault,
		StopPropagation: stopPropagation,
	})
}

func (ds *DiffSet) RemoveEventListener(t *node, name string) {
	ds.selectCurNode(t)
	ds.addItem(ditRemoveEventListener, name)
}

func (ds *DiffSet) SetRootItem(t *node, id string) {
	ds.selectCurNode(t)
	ds.addItem(ditSetRootItem, id)
}

type dtAddClientEventItem struct {
	ID                 string `json:"id"`
	ClientScriptPrefix string `json:"sp"`
	ScriptHandlerName  string `json:"sh"`
}

func (ds *DiffSet) AddClientEvent(t *node, name string, eventId string, clientScriptPrefix string, scriptHandlerName string) {
	ds.selectCurNode(t)
	ds.addItemWithKey(ditAddClientEvent, name, dtAddClientEventItem{
		ID:                 eventId,
		ClientScriptPrefix: clientScriptPrefix,
		ScriptHandlerName:  scriptHandlerName,
	})
}

func (ds *DiffSet) addMountComponent(t *node, c Component) {
	ds.selectCurNode(t)
	ds.addItem(ditMountComponent, c)
}

func (ds *DiffSet) addUnmountComponent(c Component) {
	ds.addItem(ditUnmountComponent, c)
}
