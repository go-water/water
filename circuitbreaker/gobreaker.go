package circuitbreaker

import (
	"context"
	"github.com/go-water/water/endpoint"
	"github.com/sony/gobreaker"
)

func GoBreaker(cb *gobreaker.CircuitBreaker) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (any, error) {
			return cb.Execute(func() (any, error) { return next(ctx, request) })
		}
	}
}
