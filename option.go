package water

import (
	"context"
	"golang.org/x/time/rate"
	"time"
)

type ServerOption func(*Server)

func ServerErrorHandler(errorHandler ErrorHandler) ServerOption {
	return func(s *Server) { s.errorHandler = errorHandler }
}

type ServerFinalizerFunc func(ctx context.Context, err error)

func ServerFinalizer(f ...ServerFinalizerFunc) ServerOption {
	return func(s *Server) { s.finalizer = append(s.finalizer, f...) }
}

func ServerLimiter(d time.Duration, t int) ServerOption {
	return func(s *Server) {
		s.limit = rate.NewLimiter(rate.Every(d), t)
	}
}
