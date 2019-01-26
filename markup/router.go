package markup

import (
	"strings"

	"github.com/yossoy/exciton/internal/markup"
)

type routeItem struct {
	exact    bool
	path     string
	paths    []string
	child    RenderResult
	redirect string
	fallback bool
}

func (ri *routeItem) match(srcPaths []string) (bool, map[string]string) {
	var vars map[string]string
	if len(ri.paths) > len(srcPaths) {
		return false, nil
	}
	sidx := 0
	ridx := 0
	for sidx < len(ri.paths) {
		sp := srcPaths[sidx]
		rp := ri.paths[ridx]
		var vn, val string
		switch {
		case strings.HasPrefix(rp, "*"):
			// paths:		"/aa/bb/cc/dd"
			// ri.paths:	"/aa/*v/dd"
			// vl:			len(paths) - len(ri.paths) + 1 = 4 - 3 + 1 = 2
			vl := len(srcPaths) - len(ri.paths) + 1
			vn = rp[1:]
			val = strings.Join(srcPaths[sidx:sidx+vl], "/")
			sidx += vl
		case strings.HasPrefix(rp, ":"):
			vn = rp[1:]
			val = sp
			sidx++
		default:
			if sp != rp {
				return false, nil
			}
			sidx++
		}
		if vn != "" {
			if vars == nil {
				vars = make(map[string]string)
			}
			vars[vn] = val
		}
		ridx++
	}
	if ri.exact && (sidx != len(srcPaths)) {
		return false, nil
	}
	return true, vars
}

type routing struct {
	items []*routeItem
	path  string
}

func (r *routing) procSub(src string) RenderResult {
	srcPaths := splitPath(src)
	var fallbackRoute *routeItem
	for _, ri := range r.items {
		if ri.fallback {
			if fallbackRoute != nil {
				panic("already has fallback")
			}
			fallbackRoute = ri
			continue
		}
		if ok, vars := ri.match(srcPaths); ok {
			if ri.child == nil && ri.redirect != "" {
				// redirect url
				rpaths := splitPath(ri.redirect)
				for ridx, rp := range rpaths {
					if strings.HasPrefix(rp, "*") || strings.HasPrefix(rp, ":") {
						rvn := rp[1:]
						if rv, ok := vars[rvn]; ok {
							rpaths[ridx] = rv
						} else {
							// cannot redirect
							// TODO: error render result?
							return nil
						}
					}
				}
				newSrc := strings.Join(rpaths, "/")
				return r.procSub(newSrc)
			}
			if rc, ok := ri.child.(*markup.ComponentRenderResult); ok {
				// clear old redirect vars
				markups := make([]Markup, 0, len(rc.Markups)+len(vars))
				for _, m := range rc.Markups {
					if pa, ok := m.(markup.PropApplyer); !ok || !pa.IsRedirect {
						markups = append(markups, m)
					}
				}
				for k, v := range vars {
					markups = append(markups, markup.PropApplyer{Name: k, Value: v, IsRedirect: true})
				}
				rc.Markups = markups
			}
			return ri.child
		}
	}
	if fallbackRoute != nil {
		return fallbackRoute.child
	}
	return nil
}

func (r *routing) proc(b markup.Builder) RenderResult {
	src := r.path
	if src == "" {
		src = b.Route() // BrowserRouter
	}
	return r.procSub(src)
}

func BrowserRouter(routes ...*routeItem) RenderResult {
	//TODO: validate routes
	r := &routing{
		items: routes,
	}
	return markup.FuncToRenderResult(r.proc)
}

func Router(path string, routes ...*routeItem) RenderResult {
	//TODO: validate path
	//TODO: validate routes
	r := &routing{
		items: routes,
		path:  path,
	}
	return markup.FuncToRenderResult(r.proc)
}

func Route(path string, child RenderResult) *routeItem {
	return &routeItem{
		exact: false,
		path:  path,
		paths: splitPath(path),
		child: child,
	}
}

func ExactRoute(path string, child RenderResult) *routeItem {
	r := Route(path, child)
	r.exact = true
	return r
}

func RedirectRoute(from, to string) *routeItem {
	return &routeItem{
		exact:    true,
		path:     from,
		paths:    splitPath(from),
		redirect: to,
	}
}

func FallbackRoute(child RenderResult) *routeItem {
	return &routeItem{
		child:    child,
		fallback: true,
	}
}

func splitPath(path string) []string {
	if path == "" || path == "/" {
		return nil
	}
	paths := strings.Split(path, "/")
	if len(paths) > 0 && paths[0] == "" {
		paths = paths[1:]
	}
	return paths
}

func OnClickRedirectTo(path string) EventListener {
	return markup.NewClientEventListener("click", "*exciton*", "onClickRedirectTo", []interface{}{path})
}

func Link(path string, markups ...MarkupOrChild) RenderResult {
	newMarkups := make([]MarkupOrChild, 0, len(markups)+2)
	for _, m := range markups {
		switch vt := m.(type) {
		case EventListener:
			// TODO: atode
			// if vt.Name == "click" {
			// 	panic(`Link cannot has "click" event`)
			// }
		case markup.AttrApplyer:
			if vt.Name == "href" {
				panic(`Link cannot has "href" attribute`)
			}
		}
		newMarkups = append(newMarkups, m)
	}
	el := OnClickRedirectTo(path)
	newMarkups = append(newMarkups, el)
	newMarkups = append(newMarkups, markup.AttrApplyer{Name: "href", Value: "#"})
	r, err := markup.Tag("a", newMarkups)
	if err != nil {
		panic(err)
	}
	return r
}
