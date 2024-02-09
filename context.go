package water

import "net/http"

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter

	Params Params
}

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (ctx *Context) Param(key string) string {
	return ctx.Request.PathValue(key)
}
