package markup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"testing"

	"golang.org/x/net/html/atom"

	"golang.org/x/net/html"

	"github.com/stretchr/testify/assert"
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/internal/object"
)

type testComponent1 struct {
	Core
}

func (tc1 *testComponent1) Render() *RenderResult {
	return nil
}

func getComponentInstanceFuncPtr(ci ComponentInstance) *runtime.Func {
	return runtime.FuncForPC(reflect.ValueOf(ci).Pointer())
}

func TestComponent1(t *testing.T) {
	r1, err := RegisterComponent((*testComponent1)(nil))
	assert.NoError(t, err)
	assert.NotNil(t, r1)
	r2, err := RegisterComponent((*testComponent1)(nil))
	assert.Error(t, err)
	assert.Nil(t, r2)
}

type testErrorCompoent1 struct {
}

func (tec1 testErrorCompoent1) Context() *Core                          { return nil }
func (tec1 testErrorCompoent1) Render() *RenderResult                   { return nil }
func (tec1 testErrorCompoent1) Key() interface{}                        { return nil }
func (tec1 testErrorCompoent1) Builder() *Builder                       { return nil }
func (tec1 testErrorCompoent1) Classes(classes ...string) MarkupOrChild { return nil }
func (tec1 testErrorCompoent1) ID() string                              { return "" }

type testErrorCompoent2 int

func (tec2 *testErrorCompoent2) Context() *Core                          { return nil }
func (tec2 *testErrorCompoent2) Render() *RenderResult                   { return nil }
func (tec2 *testErrorCompoent2) Key() interface{}                        { return nil }
func (tec2 *testErrorCompoent2) Builder() *Builder                       { return nil }
func (tec2 *testErrorCompoent2) Classes(classes ...string) MarkupOrChild { return nil }
func (tec2 *testErrorCompoent2) ID() string                              { return "" }

func TestComponentError1(t *testing.T) {
	f := testErrorCompoent1{}
	_, err := RegisterComponent(f)
	assert.Error(t, err)

	_, err = RegisterComponent((*testErrorCompoent2)(nil))
	assert.Error(t, err)
}

type testData struct {
	keySequence        int
	parentComponent    Component
	parentMountCount   int
	parentUnmountCount int
	child1MountCount   int
	child1UnmountCount int
	child2MountCount   int
	child2UnmountCount int
}

func makeHTMLRoot() *html.Node {
	root := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Lookup([]byte("div")),
		Data:     "div",
	}
	return root
}

func applyHTML(d *DiffSet, root *html.Node) string {
	applyDiff(d, root)
	buf := bytes.NewBufferString("")
	html.Render(buf, root)
	return buf.String()
}

func diffSetToString(d *DiffSet) string {
	bb, _ := json.MarshalIndent(d, "", "  ")
	return string(bb)
}

func (td *testData) newKey() object.ObjectKey {
	k := td.keySequence
	td.keySequence++
	return fmt.Sprintf("id%d", k)
}

type testComponentProps1 struct {
	Core
	IntValue *int `exciton:"intValue"`
}

func (tcp1 *testComponentProps1) Initialize() {
	td := tcp1.Builder().UserData.(*testData)
	if td.parentComponent != nil {
		panic("already create component instance")
	}
	td.parentComponent = tcp1
}

func (tcp1 *testComponentProps1) Render() *RenderResult {
	return Text(strconv.Itoa(*tcp1.IntValue))
}

func (tcp1 *testComponentProps1) Mount() {
	td := tcp1.Builder().UserData.(*testData)
	td.parentMountCount++
}

func (tcp1 *testComponentProps1) Unmount() {
	td := tcp1.Builder().UserData.(*testData)
	td.parentMountCount++
}

func TestComponentProps(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	ci := MustRegisterComponent((*testComponentProps1)(nil))
	b := NewBuilder()
	b.UserData = td
	b.keyGenerator = td.newKey

	v := 100
	b.RenderBody(ci(
		Property("intValue", &v),
	))
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html := applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div>100</div>`, html)

	b.diffSet.reset()
	v = 200
	b.Rerender(td.parentComponent)
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div>200</div>`, html)

	unregisterComponent(ci)
}

type nestComponentChild1 struct {
	Core
	Index int `exciton:"index"`
}

func (ncc1 *nestComponentChild1) Render() *RenderResult {
	return Tag("span",
		Text(strconv.Itoa(ncc1.Index)),
	)
}

func (ncc1 *nestComponentChild1) Mount() {
	td := ncc1.Builder().UserData.(*testData)
	td.child1MountCount++
}
func (ncc1 *nestComponentChild1) Unmount() {
	td := ncc1.Builder().UserData.(*testData)
	td.child1UnmountCount++
}

type nestComponentChild2 struct {
	Core
	Index int `exciton:"index2"`
}

func (ncc1 *nestComponentChild2) Render() *RenderResult {
	return Tag("span",
		Text(strconv.Itoa(ncc1.Index)),
	)
}
func (ncc1 *nestComponentChild2) Mount() {
	td := ncc1.Builder().UserData.(*testData)
	td.child2MountCount++
}
func (ncc1 *nestComponentChild2) Unmount() {
	td := ncc1.Builder().UserData.(*testData)
	td.child2UnmountCount++
}

type nestComponentParent1 struct {
	Core
	Var *int `exciton:"var"`
}

func (ncp1 *nestComponentParent1) Render() *RenderResult {
	var mkups []MarkupOrChild
	if *ncp1.Var == 0 {
		mkups = append(mkups, NestComponentChild1(Property("index", 0)))
		mkups = append(mkups, Data("testData", "foo"))
		mkups = append(mkups, Classes("aa"))
	} else {
		mkups = append(mkups, NestComponentChild2(Property("index2", 1)))
		mkups = append(mkups, Classes("aa", "bb"))
	}
	return Tag("div", mkups...)
}

func (ncp1 *nestComponentParent1) Initialize() {
	td := ncp1.Builder().UserData.(*testData)
	if td.parentComponent != nil {
		panic("already create component instance")
	}
	td.parentComponent = ncp1
}
func (ncc1 *nestComponentParent1) Mount() {
	td := ncc1.Builder().UserData.(*testData)
	td.parentMountCount++
}
func (ncc1 *nestComponentParent1) Unmount() {
	td := ncc1.Builder().UserData.(*testData)
	td.parentUnmountCount++
}

var NestComponentChild1 = MustRegisterComponent((*nestComponentChild1)(nil))
var NestComponentChild2 = MustRegisterComponent((*nestComponentChild2)(nil))
var NestComponentParent1 = MustRegisterComponent((*nestComponentParent1)(nil))

func TestComponentNest1(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	b := NewBuilder()
	b.UserData = td
	b.keyGenerator = td.newKey

	v := 0
	b.RenderBody(NestComponentParent1(
		Property("var", &v),
	))
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 1, td.child1MountCount)
	assert.Equal(t, 0, td.child1UnmountCount)
	assert.Equal(t, 0, td.child2MountCount)
	assert.Equal(t, 0, td.child2UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html := applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><div _uuid="id0" data-test-data="foo" class="aa"><span _uuid="id1">0</span></div></div>`, html)

	b.diffSet.reset()
	v = 1
	b.Rerender(td.parentComponent)
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 1, td.child1MountCount)
	assert.Equal(t, 1, td.child1UnmountCount)
	assert.Equal(t, 1, td.child2MountCount)
	assert.Equal(t, 0, td.child2UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><div _uuid="id0" class="aa bb"><span _uuid="id2">1</span></div></div>`, html)

	b.diffSet.reset()
	v = 0
	b.Rerender()
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 2, td.child1MountCount)
	assert.Equal(t, 1, td.child1UnmountCount)
	assert.Equal(t, 1, td.child2MountCount)
	assert.Equal(t, 1, td.child2UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><div _uuid="id0" class="aa" data-test-data="foo"><span _uuid="id3">0</span></div></div>`, html)
}

func TestComponentNestAsync1(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	var bp **Builder
	ch := make(chan string, 1)
	b := NewAsyncBuilder(
		func() {
			(*bp).ProcRequestAnimationFrame()
		}, func(ds *DiffSet) {
			t.Log(diffSetToString(ds))
			html := applyHTML(ds, root)
			ch <- html
		},
	)
	bp = &b
	b.UserData = td
	b.keyGenerator = td.newKey

	v := 0
	b.RenderBody(NestComponentParent1(
		Property("var", &v),
	))
	html := <-ch
	t.Log(html)
	assert.Equal(t, `<div><div _uuid="id0" data-test-data="foo" class="aa"><span _uuid="id1">0</span></div></div>`, html)

	b.diffSet.reset()
	v = 1
	b.Rerender(td.parentComponent)
	html = <-ch
	t.Log(html)
	assert.Equal(t, `<div><div _uuid="id0" class="aa bb"><span _uuid="id2">1</span></div></div>`, html)

	b.diffSet.reset()
	v = 0
	b.Rerender()
	html = <-ch
	t.Log(html)
	assert.Equal(t, `<div><div _uuid="id0" class="aa" data-test-data="foo"><span _uuid="id3">0</span></div></div>`, html)
}

type nestComponentParent2 struct {
	Core
	Var *int `exciton:"var"`
}

func (ncp1 *nestComponentParent2) Render() *RenderResult {
	if *ncp1.Var == 0 {
		return NestComponentChild1(Property("index", 0))
	}
	return NestComponentChild2(Property("index2", 1))
}

func (ncp1 *nestComponentParent2) Initialize() {
	td := ncp1.Builder().UserData.(*testData)
	if td.parentComponent != nil {
		panic("already create component instance")
	}
	td.parentComponent = ncp1
}

func (ncp1 *nestComponentParent2) Mount() {
	td := ncp1.Builder().UserData.(*testData)
	td.parentMountCount++
}
func (ncp1 *nestComponentParent2) Unmount() {
	td := ncp1.Builder().UserData.(*testData)
	td.parentUnmountCount++
}

var NestComponentParent2 = MustRegisterComponent((*nestComponentParent2)(nil))

func TestComponentNest2(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	b := NewBuilder()
	b.UserData = td
	b.keyGenerator = td.newKey

	v := 0
	b.RenderBody(NestComponentParent2(
		Property("var", &v),
	))
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 1, td.child1MountCount)
	assert.Equal(t, 0, td.child1UnmountCount)
	assert.Equal(t, 0, td.child2MountCount)
	assert.Equal(t, 0, td.child2UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html := applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><span _uuid="id0">0</span></div>`, html)

	b.diffSet.reset()
	v = 1
	b.Rerender(td.parentComponent)
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 1, td.child1MountCount)
	assert.Equal(t, 1, td.child1UnmountCount)
	assert.Equal(t, 1, td.child2MountCount)
	assert.Equal(t, 0, td.child2UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><span _uuid="id1">1</span></div>`, html)

	b.diffSet.reset()
	v = 0
	b.Rerender()
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 2, td.child1MountCount)
	assert.Equal(t, 1, td.child1UnmountCount)
	assert.Equal(t, 1, td.child2MountCount)
	assert.Equal(t, 1, td.child2UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><span _uuid="id2">0</span></div>`, html)
}

type nestComponentParent3 struct {
	Core
	Var *int `exciton:"var"`
}

func (ncp3 *nestComponentParent3) Render() *RenderResult {
	if *ncp3.Var == 0 {
		return NestComponentChild1(Property("index", 0))
	}
	return NestComponentChild1(Property("index", 1))
}

func (ncp3 *nestComponentParent3) Initialize() {
	td := ncp3.Builder().UserData.(*testData)
	if td.parentComponent != nil {
		panic("already create component instance")
	}
	td.parentComponent = ncp3
}
func (ncp3 *nestComponentParent3) Mount() {
	td := ncp3.Builder().UserData.(*testData)
	td.parentMountCount++
}
func (ncp3 *nestComponentParent3) Unmount() {
	td := ncp3.Builder().UserData.(*testData)
	td.parentUnmountCount++
}

var NestComponentParent3 = MustRegisterComponent((*nestComponentParent3)(nil))

func TestComponentNest3(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	b := NewBuilder()
	b.UserData = td
	b.keyGenerator = td.newKey

	v := 0
	b.RenderBody(NestComponentParent3(
		Property("var", &v),
	))
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 1, td.child1MountCount)
	assert.Equal(t, 0, td.child1UnmountCount)
	assert.Equal(t, 0, td.child2MountCount)
	assert.Equal(t, 0, td.child2UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html := applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><span _uuid="id0">0</span></div>`, html)

	b.diffSet.reset()
	v = 1
	b.Rerender(td.parentComponent)
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 1, td.child1MountCount)
	assert.Equal(t, 0, td.child1UnmountCount)
	assert.Equal(t, 0, td.child2MountCount)
	assert.Equal(t, 0, td.child2UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><span _uuid="id0">1</span></div>`, html)

	b.diffSet.reset()
	v = 0
	b.Rerender()
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 1, td.child1MountCount)
	assert.Equal(t, 0, td.child1UnmountCount)
	assert.Equal(t, 0, td.child2MountCount)
	assert.Equal(t, 0, td.child2UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><span _uuid="id0">0</span></div>`, html)
}

type nestComponentParent4 struct {
	Core
	Var *int `exciton:"var"`
}

func (ncp4 *nestComponentParent4) Render() *RenderResult {
	if *ncp4.Var == 0 {
		return NestComponentChild1(Property("index", 0))
	}
	return Text("1")
}

func (ncp4 *nestComponentParent4) Initialize() {
	td := ncp4.Builder().UserData.(*testData)
	if td.parentComponent != nil {
		panic("already create component instance")
	}
	td.parentComponent = ncp4
}
func (ncp4 *nestComponentParent4) Mount() {
	td := ncp4.Builder().UserData.(*testData)
	td.parentMountCount++
}
func (ncp4 *nestComponentParent4) Unmount() {
	td := ncp4.Builder().UserData.(*testData)
	td.parentUnmountCount++
}

var NestComponentParent4 = MustRegisterComponent((*nestComponentParent4)(nil))

func TestComponentNest4(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	b := NewBuilder()
	b.UserData = td
	b.keyGenerator = td.newKey

	v := 0
	b.RenderBody(NestComponentParent4(
		Property("var", &v),
	))
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 1, td.child1MountCount)
	assert.Equal(t, 0, td.child1UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html := applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><span _uuid="id0">0</span></div>`, html)

	b.diffSet.reset()
	v = 1
	b.Rerender(td.parentComponent)
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 1, td.child1MountCount)
	assert.Equal(t, 1, td.child1UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div>1</div>`, html)

	b.diffSet.reset()
	v = 0
	b.Rerender()
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 2, td.child1MountCount)
	assert.Equal(t, 1, td.child1UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><span _uuid="id1">0</span></div>`, html)
}

type childParentComponent1 struct {
	Core
	Var int `exciton:"var"`
}

func (cpc1 *childParentComponent1) Render() *RenderResult {
	m := make([]MarkupOrChild, len(cpc1.Children()))
	for i, c := range cpc1.Children() {
		m[i] = c
	}
	return Tag("div", m...)
}
func (cpc1 *childParentComponent1) Mount() {
	td := cpc1.Builder().UserData.(*testData)
	td.parentMountCount++
}
func (cpc1 *childParentComponent1) Unmount() {
	td := cpc1.Builder().UserData.(*testData)
	td.parentUnmountCount++
}

var ChildParentComponent1 = MustRegisterComponent((*childParentComponent1)(nil))

func TestComponentChild1(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	b := NewBuilder()
	b.UserData = td
	b.keyGenerator = td.newKey

	b.RenderBody(ChildParentComponent1(
		NestComponentChild1(Property("index", 0)),
		NestComponentChild1(Property("index", 1)),
	))
	assert.Equal(t, 1, td.parentMountCount)
	assert.Equal(t, 0, td.parentUnmountCount)
	assert.Equal(t, 2, td.child1MountCount)
	assert.Equal(t, 0, td.child1UnmountCount)
	t.Log(diffSetToString(b.diffSet))
	html := applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><div _uuid="id0"><span _uuid="id1">0</span><span _uuid="id2">1</span></div></div>`, html)
}

func makeTestEvent(name string, listener func()) *EventListener {
	return &EventListener{Name: name, Listener: func(le *event.Event) {
		listener()
	}}
}

type testEventComponent1 struct {
	Core
	Var *int `exciton:"var"`
}

func (c *testEventComponent1) Callback1() {
}
func (c *testEventComponent1) Callback2() {
}

func (c *testEventComponent1) Render() *RenderResult {
	var m MarkupOrChild
	if *c.Var == 0 {
		m = makeTestEvent("event1", c.Callback1)
	} else {
		m = makeTestEvent("event1", c.Callback2)
	}
	return Tag("button", m)
}

var EventComponent1 = MustRegisterComponent((*testEventComponent1)(nil))

func TestComponentEvent1(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	b := NewBuilder()
	b.UserData = td
	b.keyGenerator = td.newKey

	v := 0
	b.RenderBody(
		EventComponent1(Property("var", &v)),
	)
	t.Log(diffSetToString(b.diffSet))
	html := applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><button _uuid="id0" onevent1="evt{id0,false,false}"></button></div>`, html)

	b.diffSet.reset()
	v = 1
	b.Rerender()
	t.Log(diffSetToString(b.diffSet))
	assert.Equal(t, 0, len(b.diffSet.Items)) // empty diff
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><button _uuid="id0" onevent1="evt{id0,false,false}"></button></div>`, html)
}

type testEventComponent2 struct {
	Core
	Var *int `exciton:"var"`
}

func (c *testEventComponent2) Callback1() {
}

func (c *testEventComponent2) Render() *RenderResult {
	var m MarkupOrChild
	if *c.Var == 0 {
		m = makeTestEvent("event1", c.Callback1)
	}
	return Tag("button", m)
}

var EventComponent2 = MustRegisterComponent((*testEventComponent2)(nil))

func TestComponentEvent2(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	b := NewBuilder()
	b.UserData = td
	b.keyGenerator = td.newKey

	v := 0
	b.RenderBody(
		EventComponent2(Property("var", &v)),
	)
	t.Log(diffSetToString(b.diffSet))
	html := applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><button _uuid="id0" onevent1="evt{id0,false,false}"></button></div>`, html)

	b.diffSet.reset()
	v = 1
	b.Rerender()
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><button _uuid="id0"></button></div>`, html)

	b.diffSet.reset()
	v = 0
	b.Rerender()
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><button _uuid="id0" onevent1="evt{id0,false,false}"></button></div>`, html)
}

type keyTestChildComponent1 struct {
	Core
	Var int `exciton:"var"`
}

func (c *keyTestChildComponent1) Render() *RenderResult {
	if c.Var == 0 {
		return Text("a")
	}
	return nil
}

var KeyTestChildComponent1 = MustRegisterComponent((*keyTestChildComponent1)(nil))

type keyTestParentComponent1 struct {
	Core
	Var *int `exciton:"var"`
}

func (c *keyTestParentComponent1) Render() *RenderResult {
	if *c.Var == 0 {
		return Tag("div",
			Keyer("key1", KeyTestChildComponent1(Property("var", 0))),
			Keyer("key2", KeyTestChildComponent1(Property("var", 1))),
		)
	}
	return Tag("div",
		Keyer("key2", KeyTestChildComponent1(Property("var", 1))),
		Keyer("key1", KeyTestChildComponent1(Property("var", 0))),
	)
}

var KeyTestParentComponent1 = MustRegisterComponent((*keyTestParentComponent1)(nil))

type nonKeyTestParentComponent1 struct {
	Core
	Var *int `exciton:"var"`
}

func (c *nonKeyTestParentComponent1) Render() *RenderResult {
	if *c.Var == 0 {
		return Tag("div",
			KeyTestChildComponent1(Property("var", 0)),
			KeyTestChildComponent1(Property("var", 1)),
		)
	}
	return Tag("div",
		KeyTestChildComponent1(Property("var", 1)),
		KeyTestChildComponent1(Property("var", 0)),
	)
}

var NonKeyTestParentComponent1 = MustRegisterComponent((*nonKeyTestParentComponent1)(nil))

func TestComponentNonKey1(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	b := NewBuilder()
	b.UserData = td
	b.keyGenerator = td.newKey

	v := 0
	b.RenderBody(
		NonKeyTestParentComponent1(Property("var", &v)),
	)
	t.Log(diffSetToString(b.diffSet))
	html := applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><div _uuid="id0">a<noscript _uuid="id1"></noscript></div></div>`, html)

	b.diffSet.reset()
	v = 1
	b.Rerender()
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	// non keyed node recreate in renderer
	assert.Equal(t, `<div><div _uuid="id0"><noscript _uuid="id2"></noscript>a</div></div>`, html)

	b.diffSet.reset()
	v = 0
	b.Rerender()
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	// non keyed node recreate in renderer
	assert.Equal(t, `<div><div _uuid="id0">a<noscript _uuid="id3"></noscript></div></div>`, html)
}

func TestComponentKey1(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	b := NewBuilder()
	b.UserData = td
	b.keyGenerator = td.newKey

	v := 0
	b.RenderBody(
		KeyTestParentComponent1(Property("var", &v)),
	)
	t.Log(diffSetToString(b.diffSet))
	html := applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><div _uuid="id0">a<noscript _uuid="id1"></noscript></div></div>`, html)

	b.diffSet.reset()
	v = 1
	b.Rerender()
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	// keyed node keep id(node) in rendered
	assert.Equal(t, `<div><div _uuid="id0"><noscript _uuid="id1"></noscript>a</div></div>`, html)

	b.diffSet.reset()
	v = 0
	b.Rerender()
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	// keyed node keep id(node) in rendered
	assert.Equal(t, `<div><div _uuid="id0">a<noscript _uuid="id1"></noscript></div></div>`, html)
}

type innerHTMLTestComponent1 struct {
	Core
	Var *int `exciton:"var"`
}

func (c *innerHTMLTestComponent1) Render() *RenderResult {
	if *c.Var == 0 {
		return Tag("div",
			UnsafeHTML(`<span>innerHTML</span>`),
		)
	}
	return Tag("span",
		Style("color", "red"),
		Text("nonInnerHTML"),
	)
}

var InnerHTMLTestComponent1 = MustRegisterComponent((*innerHTMLTestComponent1)(nil))

func TestComponentInnerHTML1(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	b := NewBuilder()
	b.UserData = td
	b.keyGenerator = td.newKey

	v := 0
	b.RenderBody(
		InnerHTMLTestComponent1(Property("var", &v)),
	)
	t.Log(diffSetToString(b.diffSet))
	html := applyHTML(b.diffSet, root)
	t.Log(html)
	assert.Equal(t, `<div><div _uuid="id0"><span>innerHTML</span></div></div>`, html)

	b.diffSet.reset()
	v = 1
	b.Rerender()
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	// non keyed node recreate in renderer
	assert.Equal(t, `<div><span _uuid="id1" style="{&#34;color&#34;:&#34;red&#34;}">nonInnerHTML</span></div>`, html)

	b.diffSet.reset()
	v = 0
	b.Rerender()
	t.Log(diffSetToString(b.diffSet))
	html = applyHTML(b.diffSet, root)
	t.Log(html)
	// non keyed node recreate in renderer
	assert.Equal(t, `<div><div _uuid="id2"><span>innerHTML</span></div></div>`, html)
}
