package water

import (
	"context"
	"go.uber.org/zap"
)

type Handler interface {
	ServerWater(ctx context.Context, req any) (any, error)
	GetLogger() *zap.Logger
}

type Server struct {
	e            Endpoint
	finalizer    []ServerFinalizerFunc
	errorHandler ErrorHandler
	l            *zap.Logger
}

func NewHandler(srv Service, options ...ServerOption) Handler {
	s := &Server{
		e: srv.Endpoint(),
	}

	for _, option := range options {
		option(s)
	}

	handler := NewLogErrorHandler(NewLogger(), srv.Name(srv))
	srv.SetLogger(handler.l)
	s.l = handler.l
	s.errorHandler = handler

	return s
}

func (s *Server) ServerWater(ctx context.Context, req any) (resp any, err error) {
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

func (s *Server) GetLogger() *zap.Logger {
	return s.l
}
