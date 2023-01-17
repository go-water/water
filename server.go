package water

import (
	"context"
)

type Handler interface {
	ServerWater(ctx context.Context, req interface{}) (interface{}, error)
}

type Server struct {
	c            *Config
	e            Endpoint
	finalizer    []ServerFinalizerFunc
	errorHandler ErrorHandler
}

func NewServer(e Endpoint, n string, options ...ServerOption) *Server {
	s := &Server{
		e: e,
		c: new(Config),
	}

	for _, option := range options {
		option(s)
	}

	s.errorHandler = NewLogErrorHandler(s.c.NewLogger(), n)

	return s
}

func (s Server) ServerWater(ctx context.Context, req interface{}) (resp interface{}, err error) {
	if len(s.finalizer) > 0 {
		defer func() {
			for _, fn := range s.finalizer {
				fn(ctx, err)
			}
		}()
	}

	resp, err = s.e(ctx, req)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		return nil, err
	}

	return resp, nil
}
