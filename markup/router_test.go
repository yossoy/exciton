package markup

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yossoy/exciton/internal/object"
	"golang.org/x/net/html"
)

type testRoute1 struct {
	Core
	Base string `exciton:"base"`
	Var1 string `exciton:"var1"`
	Var2 string `exciton:"var2"`
}

func (c *testRoute1) Render() RenderResult {
	return Text(fmt.Sprintf("base[%s]var1[%s]:var2[%s]", c.Base, c.Var1, c.Var2))
}

var route1 = MustRegisterComponent((*testRoute1)(nil))

func testRouterSub(t *testing.T, b *Builder, root *html.Node, path string, expected string) {
	b.diffSet.reset()
	b.route = path
	b.Rerender()
	//t.Log(diffSetToString(b.diffSet))
	html := applyHTML(b.diffSet, root)
	//t.Log(html)
	// non keyed node recreate in renderer
	assert.Equal(t, expected, html)
}

func TestRouterBasic(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	b := NewBuilder("")
	b.UserData = td
	b.keyGenerator = func() object.ObjectKey { return "-" } // always -
	const unmatchedResult = `<div><noscript _uuid="-"></noscript></div>`

	b.RenderBody(
		BrowserRouter(
			ExactRoute("/", Text("000")),
			Route("/aaa", Text("aaa")),
			Route("/bbb", Text("bbb")),
		),
	)
	//	applyHTML(b.diffSet, root)
	applyHTML(b.diffSet, root)
	cases := []struct {
		route    string
		expected string
	}{
		{"/", `<div>000</div>`},
		{"/aaa", `<div>aaa</div>`},
		{"/aaaa", unmatchedResult},
		{"/bbb", `<div>bbb</div>`},
		{"/bbb/ccc", `<div>bbb</div>`},
		{"/ccc", unmatchedResult},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("route: %q", c.route), func(t *testing.T) { testRouterSub(t, b, root, c.route, c.expected) })
	}
}
func TestRouterWithVar(t *testing.T) {
	root := makeHTMLRoot()
	td := &testData{}
	b := NewBuilder("")
	b.UserData = td
	b.keyGenerator = func() object.ObjectKey { return "-" } // always -
	const unmatchedResult = `<div>fallbacked</div>`

	b.RenderBody(
		BrowserRouter(
			ExactRoute("/", Tag("span", Text("000"))),
			Route("/aaa/:var1/:var2", route1(Property("base", "/aaa"))),
			Route("/aaa/:var2", route1(Property("base", "/aaa(2)"))),
			Route("/bbb/*var1", route1(Property("base", "/bbb"))),
			Route("/ccc/:var1/ccc", route1(Property("base", "/ccc"))),
			RedirectRoute("/ddd", "/aaa/99"),
			RedirectRoute("/eee/*var1", "/aaa/:var1"),
			FallbackRoute(Text("fallbacked")),
		),
	)
	//	applyHTML(b.diffSet, root)
	applyHTML(b.diffSet, root)

	cases := []struct {
		route    string
		expected string
	}{
		{"/", `<div><span _uuid="-">000</span></div>`},
		{"/aaa", unmatchedResult},
		{"/aaa/0", `<div>base[/aaa(2)]var1[]:var2[0]</div>`},
		{"/aaa/1/2", `<div>base[/aaa]var1[1]:var2[2]</div>`},
		{"/aaa/3/4", `<div>base[/aaa]var1[3]:var2[4]</div>`},
		{"/aaa/5/6/7", `<div>base[/aaa]var1[5]:var2[6]</div>`},
		{"/aaaa", unmatchedResult},
		{"/bbb/1", `<div>base[/bbb]var1[1]:var2[]</div>`},
		{"/bbb/1/2/3", `<div>base[/bbb]var1[1/2/3]:var2[]</div>`},
		{"/ddd", `<div>base[/aaa(2)]var1[1/2/3]:var2[99]</div>`}, // not "base[/aaa(2)]var1[]:var2[99]", because keep previous var1 proerty value.
		{"/eee/11", `<div>base[/aaa(2)]var1[1/2/3]:var2[11]</div>`},
		{"/eee/111/222", `<div>base[/aaa]var1[111]:var2[222]</div>`},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("route: %q", c.route), func(t *testing.T) { testRouterSub(t, b, root, c.route, c.expected) })
	}
}
