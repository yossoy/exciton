package driver

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router is interface for registering routes to be matched and dispatches a handler.
type Router interface {
	PathPrefix(string) Route
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	HandleFunc(string, func(http.ResponseWriter, *http.Request)) Route
	Use(...MiddlewareFunc)
}

type router struct {
	r *mux.Router
}

func newRouter() *router {
	return &router{
		r: mux.NewRouter(),
	}
}

func (r *router) Use(mwf ...MiddlewareFunc) {
	r.r.Use(mwf...)
}

func (r *router) PathPrefix(tpl string) Route {
	return &route{
		r: r.r.PathPrefix(tpl),
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.r.ServeHTTP(w, req)
}

func (r *router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) Route {
	return &route{
		r: r.r.HandleFunc(path, f),
	}
}

// Route is interface for information to match a request and build URLs.
type Route interface {
	Handler(handler http.Handler) Route
	HandlerFunc(func(http.ResponseWriter, *http.Request)) Route
	PathPrefix(tpl string) Route
}

type route struct {
	r *mux.Route
}

func (r *route) Handler(handler http.Handler) Route {
	return &route{
		r: r.r.Handler(handler),
	}
}

func (r *route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) Route {
	return &route{
		r: r.r.HandlerFunc(f),
	}
}

func (r *route) PathPrefix(tpl string) Route {
	return &route{
		r: r.r.PathPrefix(tpl),
	}
}

func RequestVars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

type MiddlewareFunc = mux.MiddlewareFunc
