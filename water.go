package water

import (
	"github.com/go-water/water/render"
	"net/http"
	"path"
	"slices"
	"strings"
	"sync"
	"time"
)

type Water struct {
	Router
	ContextWithFallback bool
	HTMLRender          render.HTMLRender
	pool                sync.Pool
}

func (w *Water) allocateContext() *Context {
	return &Context{wt: w}
}

func (w *Water) Serve(addr string, server ...*http.Server) error {
	handler := &http.ServeMux{}
	for _, rt := range w.base.routes {
		for url, handle := range rt.routes {
			rhd := new(RouterHandler)
			rhd.wt = w
			rhd.h = handle
			handler.Handle(url, rhd)
		}
	}

	var h http.Handler = handler
	for _, middleware := range w.base.global {
		h = middleware(h)
	}

	srv := &http.Server{
		ReadHeaderTimeout: time.Second * 45,
	}
	if len(server) != 0 {
		srv = server[0]
	}
	srv.Addr, srv.Handler = addr, h

	return srv.ListenAndServe()
}

type Router struct {
	scope       string
	routes      map[string]HandlerFunc
	middlewares []Middleware
	base        *base

	basePath string
}

type base struct {
	global []HttpHandler
	routes map[string]*Router
}

func New() *Water {
	meili := &Water{
		Router: Router{
			routes: make(map[string]HandlerFunc),
			base: &base{
				global: make([]HttpHandler, 0),
				routes: make(map[string]*Router),
			},
		},
	}

	meili.base.routes[""] = &meili.Router
	meili.pool.New = func() any {
		return meili.allocateContext()
	}

	return meili
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

func (r *Router) POST(route string, handler HandlerFunc) {
	r.Method(http.MethodPost, route, handler)
}

func (r *Router) GET(route string, handler HandlerFunc) {
	r.Method(http.MethodGet, route, handler)
}

func (r *Router) HEAD(route string, handler HandlerFunc) {
	r.Method(http.MethodHead, route, handler)
}

func (r *Router) Put(route string, handler HandlerFunc) {
	r.Method(http.MethodPut, route, handler)
}

func (r *Router) Patch(route string, handler HandlerFunc) {
	r.Method(http.MethodPatch, route, handler)
}

func (r *Router) Delete(route string, handler HandlerFunc) {
	r.Method(http.MethodDelete, route, handler)
}

func (r *Router) OPTIONS(route string, handlers HandlerFunc) {
	r.Method(http.MethodOptions, route, handlers)
}

func (r *Router) Method(method, route string, handler HandlerFunc) {
	if strings.HasSuffix(route, "/") {
		route += "{$}"
	}

	router := r.base.routes[r.scope]
	router.routes[method+" "+r.scope+route] = r.withMiddlewares(handler)
}

func (r *Router) StaticFile(relativePath, filepath string) *Router {
	return r.staticFileHandler(relativePath, func(c *Context) {
		c.File(filepath)
	})
}

func (r *Router) staticFileHandler(relativePath string, handler HandlerFunc) *Router {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	r.GET(relativePath, handler)
	r.HEAD(relativePath, handler)
	return r
}

func (r *Router) Static(relativePath, root string) *Router {
	return r.StaticFS(relativePath, Dir(root, false))
}

func (r *Router) StaticFS(relativePath string, fs http.FileSystem) *Router {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := r.createStaticHandler(relativePath, fs)
	urlPattern := path.Join(relativePath, "{path...}")

	r.GET(urlPattern, handler)
	r.HEAD(urlPattern, handler)
	return r
}

func (r *Router) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := r.calculateAbsolutePath(relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {
		//if _, noListing := fs.(*onlyFilesFS); noListing {
		//	c.Writer.WriteHeader(http.StatusNotFound)
		//}

		file := c.Param("path")
		// Check if file exists and/or if we have permission to access it
		f, err := fs.Open(file)
		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			//c.handlers = r.engine.noRoute
			// Reset index
			//c.index = -1
			return
		}
		f.Close()

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (r *Router) calculateAbsolutePath(relativePath string) string {
	return joinPaths(r.basePath, relativePath)
}

func (r *Router) Use(middlewares ...Middleware) {
	slices.Reverse(middlewares)
	r.middlewares = slices.Concat(middlewares, r.middlewares)
}

func (r *Router) withMiddlewares(handler HandlerFunc) HandlerFunc {
	for _, middleware := range r.middlewares {
		handler = middleware(handler)
	}
	return handler
}
