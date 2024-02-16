package water

import (
	"context"
	"errors"
	"github.com/go-water/water/circuitbreaker"
	"github.com/go-water/water/endpoint"
	"github.com/go-water/water/ratelimit"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"log/slog"
	"reflect"
)

type adapter struct {
	e         endpoint.Endpoint
	finalizer []ServerFinalizerFunc
	l         *slog.Logger
	limit     *rate.Limiter
	breaker   *gobreaker.CircuitBreaker
}

func NewHandler(srv Service, options ...ServerOption) Handler {
	a := new(adapter)
	for _, option := range options {
		option(a)
	}

	a.e = a.endpoint(srv)
	if a.limit != nil {
		a.e = ratelimit.NewErrorLimiter(a.limit)(a.e)
	}
	if a.breaker != nil {
		a.e = circuitbreaker.GoBreaker(a.breaker)(a.e)
	}

	l := Logger.With(slog.String("name", srv.Name(srv)))
	srv.SetLogger(l)
	a.l = l

	return a
}

func (a *adapter) endpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		function, srv, ctxV, reqV, err := a.readRequest(ctx, service, req)
		if err != nil {
			return nil, err
		}

		returnValues := function.Call([]reflect.Value{srv, ctxV, reqV})
		if len(returnValues) != 2 {
			return nil, errors.New("method Handle does not return two arguments")
		}

		returnValue := returnValues[1].Interface()
		if returnValue == nil {
			return returnValues[0].Interface(), nil
		} else {
			er, ok := returnValue.(error)
			if ok {
				return nil, er
			}

			return nil, errors.New("method Handle return argument not include error type")
		}
	}
}

func (a *adapter) readRequest(ctx context.Context, service Service, req any) (function, srv, ctxV, reqV reflect.Value, err error) {
	typ := reflect.TypeOf(service)
	srv = reflect.ValueOf(service)

	method, ok := typ.MethodByName("Handle")
	if ok {
		mType := method.Type
		num := mType.NumIn()
		if num == 3 {
			ctxV = reflect.ValueOf(ctx)
			reqV = reflect.ValueOf(req)
			function = method.Func
		} else {
			err = errors.New("method Handle does not include two parameters")
		}
	} else {
		err = errors.New("method Handle not implemented")
	}

	return
}

func (a *adapter) ServerWater(ctx context.Context, req any) (resp any, err error) {
	if len(a.finalizer) > 0 {
		defer func() {
			for _, fn := range a.finalizer {
				fn(ctx, err)
			}
		}()
	}

	resp, err = a.e(ctx, req)
	if err != nil {
		a.l.Error(err.Error())
		return nil, err
	}

	return resp, nil
}

func (a *adapter) GetLogger() *slog.Logger {
	return a.l
}
