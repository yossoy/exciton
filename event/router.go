package event

import (
	"strings"

	"github.com/pkg/errors"
)

// RouteItem is interface for router's item
type RouteItem interface {
	Item() interface{}
	IsUnmatched() bool
}

// UnmatchedRouteItem is interface
type UnmatchedRouteItem interface {
	RouteItem
	PathSegments() string
}

type routeItem struct {
	substrs  []string
	item     interface{}
	children *Router
}

func (ri *routeItem) Item() interface{} {
	return ri.item
}

func (ri *routeItem) IsUnmatched() bool {
	return false
}

type unmatchedItem struct {
	substrs []string
	item    interface{}
}

func (ui *unmatchedItem) Item() interface{} {
	return ui.item
}

func (ui *unmatchedItem) IsUnmatched() bool {
	return true
}

func (ui *unmatchedItem) PathSegments() string {
	return joinPathSegments(ui.substrs)
}

func (ri *routeItem) matchPath(ss []string) bool {
	if len(ss) != len(ri.substrs) {
		return false
	}
	for idx, s := range ss {
		if s != ri.substrs[idx] {
			return false
		}
	}
	return true
}

func (ri *routeItem) match(ss []string) (bool, map[string]string) {
	if len(ss) != len(ri.substrs) {
		return false, nil
	}
	var params map[string]string
	for idx, s := range ri.substrs {
		switch {
		case len(ss[idx]) == 0:
			return false, nil
		case s[0] == ':':
			if params == nil {
				params = make(map[string]string)
			}
			params[s[1:]] = ss[idx]
		default:
			if s != ss[idx] {
				return false, nil
			}
		}
	}

	return true, params
}

// Router is route item by name
type Router struct {
	routes        []*routeItem
	unmatchedItem interface{}
}

func joinPathSegments(ss []string) string {
	return "/" + strings.Join(ss, "/")
}

func splitRouterPath(s string) ([]string, error) {
	if s == "" {
		return nil, errors.New("empty path")
	}
	ss := strings.Split(s, "/")
	if len(ss) == 0 || len(ss[0]) != 0 {
		return nil, errors.New("invalid path format: " + s)
	}
	ss = ss[1:]
	for _, sss := range ss {
		if len(sss) == 0 {
			return nil, errors.New("invalid path format: " + s)
		}
	}
	return ss, nil
}

func (rr *Router) matchCore(ss []string) (RouteItem, map[string]string, []string, error) {
	for _, r := range rr.routes {
		if len(r.substrs) <= len(ss) {
			if ok, params := r.match(ss[0:len(r.substrs)]); ok {
				if r.children != nil {
					rrr, ppp, sss, err := r.children.matchCore(ss[len(r.substrs):])
					if rrr != nil {
						if params == nil {
							params = ppp
						} else {
							for k, v := range ppp {
								params[k] = v
							}
						}
						return rrr, params, append(r.substrs, sss...), err
					}
				} else if len(r.substrs) == len(ss) {
					return r, params, r.substrs, nil
				}
			}
		}
	}
	if rr.unmatchedItem != nil {
		return &unmatchedItem{
				substrs: ss,
				item:    rr.unmatchedItem,
			},
			nil, //TODO:
			nil,
			nil
	}
	return nil, nil, nil, nil
}

// Match match
func (rr *Router) Match(s string) (RouteItem, map[string]string, string, error) {
	ss, err := splitRouterPath(s)
	if err != nil {
		return nil, nil, "", err
	}
	r, p, sss, err := rr.matchCore(ss)
	return r, p, joinPathSegments(sss), err
}

func (rr *Router) isExistPath(ss []string) bool {
	for _, r := range rr.routes {
		if len(ss) != len(r.substrs) {
			continue
		}
		matched := true
		for idx, rs := range r.substrs {
			s := ss[idx]
			if s[0] != ':' && rs[0] != ':' {
				if s != rs {
					matched = false
					break
				}
			}
		}
		if matched {
			return true
		}
	}
	return false
}

// Add add route item
func (rr *Router) Add(s string, item interface{}) error {
	ss, err := splitRouterPath(s)
	if err != nil {
		return err
	}
	if rr.isExistPath(ss) {
		return errors.New("already registerd: " + s) //TODO:
	}
	r := &routeItem{
		substrs: ss,
		item:    item,
	}
	rr.routes = append(rr.routes, r)
	return nil
}

// Delete delete route item
func (rr *Router) Delete(s string) error {
	ss, err := splitRouterPath(s)
	if err != nil {
		return err
	}
	for idx, r := range rr.routes {
		if r.matchPath(ss) {
			rr.routes = append(rr.routes[:idx], rr.routes[idx+1:]...)
			return nil
		}
	}
	return errors.New("path not found")
}

// AddRoute add sub route
func (rr *Router) AddRoute(s string, router *Router) error {
	ss, err := splitRouterPath(s)
	if err != nil {
		return err
	}
	if rr.isExistPath(ss) {
		return errors.New("already registerd: " + s) //TODO:
	}
	r := &routeItem{
		substrs:  ss,
		children: router,
	}
	rr.routes = append(rr.routes, r)
	return nil
}

// SetUnmatchedItem set results used path is unmatched
func (rr *Router) SetUnmatchedItem(item interface{}) {
	rr.unmatchedItem = item
}

// NewRouter create new router
func NewRouter() *Router {
	return &Router{}
}
