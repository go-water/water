package water

import (
	"context"
	"time"

	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type ServerOption func(h *handler)

type Filter func(ctx context.Context) error

type FilterFunc func() Filter

func ServerFilterFunc(fn FilterFunc) ServerOption {
	return func(h *handler) { h.filter = fn() }
}

type FinalizerFunc func(ctx context.Context, err error)

func ServerFinalizer(f ...FinalizerFunc) ServerOption {
	return func(h *handler) { h.finalizer = append(h.finalizer, f...) }
}

func ServerErrorLimiter(interval time.Duration, b int) ServerOption {
	return func(h *handler) {
		h.el = rate.NewLimiter(rate.Every(interval), b)
	}
}

func ServerDelayLimiter(interval time.Duration, b int) ServerOption {
	return func(h *handler) {
		h.dl = rate.NewLimiter(rate.Every(interval), b)
	}
}

func ServerBreaker(breaker *gobreaker.CircuitBreaker) ServerOption {
	return func(h *handler) {
		h.breaker = breaker
	}
}
