package water

import "net/http"

type HandlerFunc func(*Context)

type RouterHandler struct {
	wt *Water
	h  HandlerFunc
}

func (r *RouterHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := r.wt.pool.Get().(*Context)
	ctx.Writer = w
	ctx.Request = req
	ctx.wt = r.wt
	ctx.reset()

	r.h(ctx)
}
