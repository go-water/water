package water

import (
	"context"
	"github.com/sony/gobreaker"
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

func ServerLimiter(interval time.Duration, b int) ServerOption {
	return func(s *Server) {
		s.limit = rate.NewLimiter(rate.Every(interval), b)
	}
}

func ServerBreaker(breaker *gobreaker.CircuitBreaker) ServerOption {
	return func(s *Server) {
		s.breaker = breaker
	}
}
