package water

import (
	"go.uber.org/zap"
	"net/http"
	"slices"
	"strings"
	"time"
)

type Service interface {
	Name(srv Service) string
	SetLogger(l *zap.Logger)
}

type Router struct {
	scope       string
	routes      map[string]HandlerFunc
	middlewares []func(HandlerFunc) HandlerFunc
	base        *base
}

type base struct {
	global []func(HandlerFunc) HandlerFunc
	routes map[string]*Router
}

type HandlerFunc func(*Context)

func NewRouter() *Router {
	rt := &Router{
		routes: make(map[string]HandlerFunc),
		base: &base{
			global: make([]func(HandlerFunc) HandlerFunc, 0),
			routes: make(map[string]*Router),
		},
	}
	rt.base.routes[""] = rt
	return rt
}

func (r *Router) Group(prefix string) *Router {
	newScope := r.scope + prefix
	newRouter := &Router{
		scope:       newScope,
		routes:      make(map[string]HandlerFunc),
		middlewares: r.middlewares,
		base:        r.base,
	}

	r.base.routes[newScope] = newRouter
	return newRouter
}

func (r *Router) Post(route string, handler HandlerFunc) {
	r.Method(http.MethodPost, route, handler)
}

func (r *Router) Get(route string, handler HandlerFunc) {
	r.Method(http.MethodGet, route, handler)
}

func (r *Router) Method(method, route string, handler HandlerFunc) {
	if strings.HasSuffix(route, "/") {
		route += "{$}"
	}

	rt := r.base.routes[r.scope]
	rt.routes[method+" "+r.scope+route] = r.withMiddlewares(handler)
}

func (r *Router) Use(middlewares ...func(HandlerFunc) HandlerFunc) {
	slices.Reverse(middlewares)
	r.middlewares = slices.Concat(middlewares, r.middlewares)
}

func (handle HandlerFunc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := new(Context)
	ctx.Writer = w
	ctx.Request = req
	handle(ctx)
}

func (r *Router) Serve(addr string, server ...*http.Server) error {
	handler := &http.ServeMux{}
	for _, rt := range r.base.routes {
		for path, handle := range rt.routes {
			handler.Handle(path, handle)
		}
	}

	srv := &http.Server{
		ReadHeaderTimeout: time.Second * 45,
	}
	if len(server) != 0 {
		srv = server[0]
	}
	srv.Addr, srv.Handler = addr, handler

	return srv.ListenAndServe()
}

func (r *Router) withMiddlewares(handler HandlerFunc) HandlerFunc {
	for _, middleware := range r.middlewares {
		handler = middleware(handler)
	}
	return handler
}
