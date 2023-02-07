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

func NewServer(srv Service, options ...ServerOption) *Server {
	s := &Server{
		e: srv.Endpoint(),
		c: new(Config),
	}

	for _, option := range options {
		option(s)
	}

	handler := NewLogErrorHandler(s.c.NewLogger(), srv.Name())
	srv.SetLogger(handler.logger)
	s.errorHandler = handler

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
