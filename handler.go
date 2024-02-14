package water

import "net/http"

type HandlerFunc func(*Context)

type MeiliHandler struct {
	m *Meili
	h HandlerFunc
}

func (mh *MeiliHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := mh.m.pool.Get().(*Context)
	ctx.Writer = w
	ctx.Request = req
	ctx.meili = mh.m
	ctx.reset()

	mh.h(ctx)
}
