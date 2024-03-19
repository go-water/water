package ratelimit

import (
	"context"
	"github.com/go-water/water/consterr"
	"github.com/go-water/water/endpoint"
)

type Allower interface {
	Allow() bool
}

func NewErrorLimiter(limit Allower) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (any, error) {
			if !limit.Allow() {
				return nil, consterr.ErrLimited
			}
			return next(ctx, request)
		}
	}
}

type Waiter interface {
	Wait(ctx context.Context) error
}

func NewDelayingLimiter(limit Waiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (any, error) {
			if err := limit.Wait(ctx); err != nil {
				return nil, err
			}
			return next(ctx, request)
		}
	}
}
