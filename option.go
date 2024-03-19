package water

import (
	"context"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"time"
)

type ServerOption func(h *handler)

type ServerFinalizerFunc func(ctx context.Context, err error)

func ServerFinalizer(f ...ServerFinalizerFunc) ServerOption {
	return func(h *handler) { h.finalizer = append(h.finalizer, f...) }
}

func ServerLimiter(interval time.Duration, b int) ServerOption {
	return func(h *handler) {
		h.limit = rate.NewLimiter(rate.Every(interval), b)
	}
}

func ServerDelayLimiter(interval time.Duration, b int) ServerOption {
	return func(h *handler) {
		h.delay = rate.NewLimiter(rate.Every(interval), b)
	}
}

func ServerBreaker(breaker *gobreaker.CircuitBreaker) ServerOption {
	return func(h *handler) {
		h.breaker = breaker
	}
}
