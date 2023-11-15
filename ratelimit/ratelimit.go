package ratelimit

import (
	"context"
	"github.com/go-water/water/consterr"
	"github.com/go-water/water/endpoint"
)

type Allowing interface {
	Allow() bool
}

func NewErrorLimiter(limit Allowing) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if !limit.Allow() {
				return nil, consterr.ErrLimited
			}
			return next(ctx, request)
		}
	}
}

type Waiting interface {
	Wait(ctx context.Context) error
}

func NewDelayingLimiter(limit Waiting) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if err := limit.Wait(ctx); err != nil {
				return nil, err
			}
			return next(ctx, request)
		}
	}
}

type AllowingFunc func() bool

func (f AllowingFunc) Allow() bool {
	return f()
}

type WaitingFunc func(ctx context.Context) error

func (f WaitingFunc) Wait(ctx context.Context) error {
	return f(ctx)
}
