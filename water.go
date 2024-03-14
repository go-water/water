package water

import (
	"fmt"
	"github.com/go-water/water/render"
	"net"
	"net/http"
	"path"
	"slices"
	"strings"
	"sync"
	"time"
)

const defaultMultipartMemory = 32 << 20 // 32 MB

type Water struct {
	Router
	ContextWithFallback bool
	HTMLRender          render.HTMLRender
	pool                sync.Pool
	TrustedPlatform     string
	RemoteIPHeaders     []string

	MaxMultipartMemory int64
}

func New() *Water {
	w := &Water{
		Router: Router{
			routes: make(map[string]HandlerFunc),
			base: &base{
				global: make([]HttpHandler, 0),
				routes: make(map[string]*Router),
			},
		},
		RemoteIPHeaders:    []string{"X-Forwarded-For", "X-Real-IP"},
		MaxMultipartMemory: defaultMultipartMemory,
	}

	w.base.routes[""] = &w.Router
	w.pool.New = func() any {
		return w.allocateContext()
	}

	return w
}

func (w *Water) Run(addr string, server ...*http.Server) error {
	mux := &http.ServeMux{}
	for _, rt := range w.base.routes {
		for url, handle := range rt.routes {
			rhd := new(RouterHandler)
			rhd.wt = w
			rhd.h = handle
			mux.Handle(url, rhd)
		}
	}

	var h http.Handler = mux
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
	w.print(addr)
	return srv.ListenAndServe()
}

func (w *Water) allocateContext() *Context {
	return &Context{wt: w}
}

func (w *Water) validateHeader(header string) (clientIP string, valid bool) {
	if header == "" {
		return "", false
	}

	items := strings.Split(header, ",")
	for i := len(items) - 1; i >= 0; i-- {
		ipStr := strings.TrimSpace(items[i])
		ip := net.ParseIP(ipStr)
		if ip == nil {
			break
		}

		if i == 0 {
			return ipStr, true
		}
	}

	return "", false
}

func (w *Water) print(addr string) {
	fmt.Println(" _   _      _ _         __        __         _     _ _ ")
	fmt.Println("| | | | ___| | | ___    \\ \\      / /__  _ __| | __| | |")
	fmt.Println("| |_| |/ _ \\ | |/ _ \\    \\ \\ /\\ / / _ \\| '__| |/ _` | |")
	fmt.Println("|  _  |  __/ | | (_) |    \\ V  V / (_) | |  | | (_| |_|")
	fmt.Println("|_| |_|\\___|_|_|\\___( )    \\_/\\_/ \\___/|_|  |_|\\__,_(_)")
	fmt.Println("                    |/")
	fmt.Println(fmt.Sprintf("Listening and serving HTTP on %s", addr))
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
	hdl := r.createStaticHandler(relativePath, fs)
	urlPattern := path.Join(relativePath, "{path...}")

	r.GET(urlPattern, hdl)
	r.HEAD(urlPattern, hdl)
	return r
}

func (r *Router) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := r.calculateAbsolutePath(relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {
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
