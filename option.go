package water

import (
	"context"
	"time"

	"github.com/go-water/water/ratelimit"
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
		h.el = ratelimit.NewErrorLimiter(rate.NewLimiter(rate.Every(interval), b))
	}
}

func ServerDelayLimiter(interval time.Duration, b int) ServerOption {
	return func(h *handler) {
		h.dl = ratelimit.NewDelayingLimiter(rate.NewLimiter(rate.Every(interval), b))
	}
}

func ServerUserErrorLimiter(interval time.Duration, b int, fn func(ctx context.Context) string) ServerOption {
	return func(h *handler) {
		h.eul = ratelimit.NewUserBasedLimiter(interval, b).UserErrorLimiter(fn)
	}
}

func ServerUserDelayLimiter(interval time.Duration, b int, fn func(ctx context.Context) string) ServerOption {
	return func(h *handler) {
		h.dul = ratelimit.NewUserBasedLimiter(interval, b).UserDelayingLimiter(fn)
	}
}

func ServerIPErrorLimiter(interval time.Duration, b int, fn func(ctx context.Context) string) ServerOption {
	return func(h *handler) {
		h.eil = ratelimit.NewIPBasedLimiter(interval, b).IPErrorLimiter(fn)
	}
}

func ServerIPDelayLimiter(interval time.Duration, b int, fn func(ctx context.Context) string) ServerOption {
	return func(h *handler) {
		h.dil = ratelimit.NewIPBasedLimiter(interval, b).IPDelayingLimiter(fn)
	}
}

func ServerBreaker(breaker *gobreaker.CircuitBreaker) ServerOption {
	return func(h *handler) {
		h.breaker = breaker
	}
}
