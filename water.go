package water

import (
	"github.com/go-water/water/render"
	"go.uber.org/zap"
	"net/http"
	"path"
	"slices"
	"strings"
	"sync"
	"time"
)

type Service interface {
	Name(srv Service) string
	SetLogger(l *zap.Logger)
}

type Router struct {
	HTMLRender  render.HTMLRender
	scope       string
	routes      map[string]HandlerFunc
	middlewares []func(HandlerFunc) HandlerFunc
	base        *base
	pool        sync.Pool

	basePath            string
	ContextWithFallback bool
}

type base struct {
	global []func(HandlerFunc) HandlerFunc
	routes map[string]*Router
}

var rt *Router

type HandlerFunc func(*Context)

func NewRouter() *Router {
	rt = &Router{
		routes: make(map[string]HandlerFunc),
		base: &base{
			global: make([]func(HandlerFunc) HandlerFunc, 0),
			routes: make(map[string]*Router),
		},
	}
	rt.base.routes[""] = rt
	rt.pool.New = func() any {
		return rt.allocateContext()
	}
	return rt
}

func (r *Router) allocateContext() *Context {
	return &Context{Router: r}
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

func (r *Router) Method(method, route string, handler HandlerFunc) {
	if strings.HasSuffix(route, "/") {
		route += "{$}"
	}

	rt := r.base.routes[r.scope]
	rt.routes[method+" "+r.scope+route] = r.withMiddlewares(handler)
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

// router.Static("/static", "/var/www")
func (r *Router) Static(relativePath, root string) *Router {
	return r.StaticFS(relativePath, Dir(root, false))
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
// Gin by default uses: gin.Dir()
func (r *Router) StaticFS(relativePath string, fs http.FileSystem) *Router {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := r.createStaticHandler(relativePath, fs)
	urlPattern := path.Join(relativePath, "{path...}")

	// Register GET and HEAD handlers
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

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func (r *Router) Use(middlewares ...func(HandlerFunc) HandlerFunc) {
	slices.Reverse(middlewares)
	r.middlewares = slices.Concat(middlewares, r.middlewares)
}

func (handle HandlerFunc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := new(Context)
	ctx.Writer = w
	ctx.Request = req
	ctx.Router = rt

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
