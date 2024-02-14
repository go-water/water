package water

type Middleware func(HandlerFunc) HandlerFunc
