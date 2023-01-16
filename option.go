package water

import "github.com/gin-gonic/gin"

type ServerOption func(*Server)

func ServerErrorHandler(errorHandler ErrorHandler) ServerOption {
	return func(s *Server) { s.errorHandler = errorHandler }
}

type ServerFinalizerFunc func(ctx *gin.Context, err error)

func ServerFinalizer(f ...ServerFinalizerFunc) ServerOption {
	return func(s *Server) { s.finalizer = append(s.finalizer, f...) }
}

func ServerConfig(c *Config) ServerOption {
	return func(s *Server) { s.c = c }
}
