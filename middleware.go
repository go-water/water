package water

import "net/http"

type Middleware func(HandlerFunc) HandlerFunc
type HttpHandler func(http.Handler) http.Handler
