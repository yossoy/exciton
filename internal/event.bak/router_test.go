package event

import (
	"reflect"
	"testing"
)

func TestRouterSimple1(t *testing.T) {
	r := NewRouter()
	var item1 interface{} = 1
	var item2 interface{} = 2
	var item3 interface{} = 3
	var item4 interface{} = 4
	var item5 interface{} = 5
	var item6 interface{} = 6
	sources := []struct {
		path    string
		isError bool
		item    interface{}
	}{
		{"/test", false, item1},
		{"", true, nil},
		{"/", true, nil},
		{"a/b", true, nil},
		{"/test", true, nil},
		{"/param1/:param", false, item2},
	}
	for n, s := range sources {
		err := r.Add(s.path, s.item)
		if s.isError {
			if err == nil {
				t.Errorf("[%d][%s] router add: invalid succeeded\n", n, s.path)
			}
		} else {
			if err != nil {
				t.Errorf("[%d][%s] router add: failed => %#v\n", n, s.path, err)
			}
		}
	}

	r.AddRoute("/parent", func() *Router {
		rc := NewRouter()
		rc.Add("/child1", item3)
		rc.Add("/child2/:param", item4)
		return rc
	}())
	r.AddRoute("/p2/:param1", func() *Router {
		rc := NewRouter()
		rc.Add("/child1", item5)
		rc.Add("/child2/:param2", item6)
		return rc
	}())

	cases := []struct {
		path    string
		isError bool
		item    interface{}
		params  map[string]string
	}{
		{"/test", false, item1, nil},
		{"/foo", false, nil, nil},
		{"", true, nil, nil},
		{"/", true, nil, nil},
		{"a/b", true, nil, nil},
		{"/test/a", false, nil, nil},
		{"/param1/foo", false, item2, map[string]string{"param": "foo"}},
		{"/parent/child1", false, item3, nil},
		{"/parent/child2/var", false, item4, map[string]string{"param": "var"}},
		{"/p2/foo/child1", false, item5, map[string]string{"param1": "foo"}},
		{"/p2/foo/child2/bar", false, item6, map[string]string{"param1": "foo", "param2": "bar"}},
	}

	for n, c := range cases {
		m, p, rn, err := r.Match(c.path)
		if c.isError {
			if err == nil {
				t.Errorf("[%d][%s] error not found:", n, c.path)
			}
		} else {
			if err != nil {
				t.Errorf("[%d][%s] invalid error: %q", n, c.path, err)
			}
		}
		if c.item != nil {
			if m.Item() != c.item {
				t.Errorf("[%d][%s] result not match: %v, %v", n, c.path, c.item, m)
			}
		} else {
			if m != nil {
				t.Errorf("[%d][%s] invalid match result: %v", n, c.path, m)
			}
		}
		if !reflect.DeepEqual(c.params, p) {
			t.Errorf("[%d][%s] result params not match: %v, %v\n", n, c.path, c.params, p)
		}
		if err == nil {
			t.Logf("[%d][%s] route name :%q\n", n, c.path, rn)
		}
	}

	err := r.Delete("/test")
	if err != nil {
		t.Errorf("remove /test is failed: %v", err)
	}
	m, _, _, err := r.Match("/test")
	if err != nil {
		t.Errorf("removed /test match is errord: %v\n", err)
	} else if m != nil {
		t.Errorf("removed /test match one: %v\n", m)
	}
	err = r.Delete("/foo")
	if err == nil {
		t.Errorf("removed /foo is not failed")
	}
}

func TestRouterUnmatched1(t *testing.T) {
	r := NewRouter()
	var item1 interface{} = 1
	r2 := NewRouter()
	var item2 interface{} = 2
	var item3 interface{} = 3
	r2.SetUnmatchedItem(item2)
	r2.Add("/bar", item3)
	r.Add("/r1/:foo", item1)
	r.AddRoute("/r2/:id", r2)

	cases := []struct {
		path         string
		isError      bool
		isFallback   bool
		item         interface{}
		params       map[string]string
		fallbackPath string
	}{
		{"/test", false, false, nil, nil, ""},
		{"/r1/foo", false, false, item1, map[string]string{"foo": "foo"}, ""},
		{"/r1/hoge/hoge", false, false, nil, nil, ""},
		{"/r2/foo/bar", false, false, item3, map[string]string{"id": "foo"}, ""},
		{"/r2/foo/baz", false, true, item2, map[string]string{"id": "foo"}, "/baz"},
	}
	for n, c := range cases {
		m, p, rn, err := r.Match(c.path)
		if c.isError {
			if err == nil {
				t.Errorf("[%d][%s] error not found:", n, c.path)
			}
		} else {
			if err != nil {
				t.Errorf("[%d][%s] invalid error: %q", n, c.path, err)
			}
		}
		if c.item != nil {
			if m.Item() != c.item {
				t.Errorf("[%d][%s] result not match: %v, %v", n, c.path, c.item, m)
			}
			if c.isFallback {
				if !m.IsUnmatched() {
					t.Errorf("[%d][%s] invalid match type: %v", n, c.path, m)
				} else {
					um, ok := m.(UnmatchedRouteItem)
					if !ok {
						t.Errorf("[%d][%s] invalid match item type: %v", n, c.path, m)
					} else {
						if um.PathSegments() != c.fallbackPath {
							t.Errorf("[%d][%s] invalid unmatched path: %q vs %q", n, c.params, c.fallbackPath, um.PathSegments())
						}
					}
				}
			} else {
				if m.IsUnmatched() {
					t.Errorf("[%d][%s] invalid match type: %v", n, c.path, m)
				}
			}
		} else {
			if m != nil {
				t.Errorf("[%d][%s] invalid match result: %v", n, c.path, m)
			}
		}
		if !reflect.DeepEqual(c.params, p) {
			t.Errorf("[%d][%s] result params not match: %v, %v\n", n, c.path, c.params, p)
		}
		if err == nil {
			t.Logf("[%d][%s] route name :%q(%v)\n", n, c.path, rn, m)
		}
	}
}
