package water

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-water/water/logger"
)

type Handler interface {
	ServeGin(ctx *gin.Context, req interface{}) (interface{}, error)
}

type Server struct {
	e            Endpoint
	finalizer    []ServerFinalizerFunc
	errorHandler ErrorHandler
}

func NewServer(e Endpoint) *Server {
	zLog := logger.Config{}
	s := &Server{
		e:            e,
		errorHandler: NewLogErrorHandler(zLog.NewLogger()),
	}

	return s
}

func (s Server) ServeGin(ctx *gin.Context, req interface{}) (interface{}, error) {
	resp, err := s.e(ctx, req)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		return nil, err
	}

	return resp, nil
}

type ServerOption func(*Server)

func ServerErrorHandler(errorHandler ErrorHandler) ServerOption {
	return func(s *Server) { s.errorHandler = errorHandler }
}

type ServerFinalizerFunc func(ctx context.Context, err error)

func ServerFinalizer(f ...ServerFinalizerFunc) ServerOption {
	return func(s *Server) { s.finalizer = append(s.finalizer, f...) }
}
