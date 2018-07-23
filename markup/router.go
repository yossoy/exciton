package markup

import (
	"strings"
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
			if rc, ok := ri.child.(*componentRenderResult); ok {
				// clear old redirect vars
				markups := make([]Markup, 0, len(rc.markups)+len(vars))
				for _, m := range rc.markups {
					if pa, ok := m.(propApplyer); !ok || !pa.isRedirect {
						markups = append(markups, m)
					}
				}
				for k, v := range vars {
					markups = append(markups, propApplyer{name: k, value: v, isRedirect: true})
				}
				rc.markups = markups
			}
			return ri.child
		}
	}
	if fallbackRoute != nil {
		return fallbackRoute.child
	}
	return nil
}

func (r *routing) proc(b *Builder) RenderResult {
	src := r.path
	if src == "" {
		src = b.route // BrowserRouter
	}
	return r.procSub(src)
}

func (r *routing) compare(n *node, hydrating bool) bool {
	return false
}

func BrowserRouter(routes ...*routeItem) RenderResult {
	//TODO: validate routes
	r := &routing{
		items: routes,
	}
	drr := &delayRenderResult{
		data:    r,
		proc:    r.proc,
		compare: r.compare,
	}
	return drr
}

func Router(path string, routes ...*routeItem) RenderResult {
	//TODO: validate path
	//TODO: validate routes
	r := &routing{
		items: routes,
		path:  path,
	}
	drr := &delayRenderResult{
		data:    r,
		proc:    r.proc,
		compare: r.compare,
	}
	return drr
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

func OnClickRedirectTo(path string) *EventListener {
	el := &EventListener{
		Name:               "click",
		clientScriptPrefix: "*exciton*",
		scriptHandlerName:  "onClickRedirectTo",
		scriptArguments:    []interface{}{path},
	}
	return el
}

func Link(path string, markups ...MarkupOrChild) RenderResult {
	newMarkups := make([]MarkupOrChild, 0, len(markups)+2)
	for _, m := range markups {
		switch vt := m.(type) {
		case *EventListener:
			if vt.Name == "click" {
				panic(`Link cannot has "click" event`)
			}
		case attrApplyer:
			if vt.name == "href" {
				panic(`Link cannot has "href" attribute`)
			}
		}
		newMarkups = append(newMarkups, m)
	}
	el := OnClickRedirectTo(path)
	newMarkups = append(newMarkups, el)
	newMarkups = append(newMarkups, attrApplyer{name: "href", value: "#"})
	r, err := tag("a", newMarkups)
	if err != nil {
		panic(err)
	}
	return r
}
