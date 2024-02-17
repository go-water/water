package water

import (
	"context"
	"errors"
	"github.com/go-water/water/circuitbreaker"
	"github.com/go-water/water/endpoint"
	"github.com/go-water/water/logger"
	"github.com/go-water/water/ratelimit"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"log/slog"
	"reflect"
)

type handler struct {
	e         endpoint.Endpoint
	finalizer []ServerFinalizerFunc
	l         *slog.Logger
	limit     *rate.Limiter
	breaker   *gobreaker.CircuitBreaker
}

func NewHandler(srv Service, options ...ServerOption) Handler {
	h := new(handler)
	for _, option := range options {
		option(h)
	}

	h.e = h.endpoint(srv)
	if h.limit != nil {
		h.e = ratelimit.NewErrorLimiter(h.limit)(h.e)
	}
	if h.breaker != nil {
		h.e = circuitbreaker.GoBreaker(h.breaker)(h.e)
	}

	l := logger.NewLogger(logger.Level, logger.AddSource).With(slog.String("name", srv.Name(srv)))
	srv.SetLogger(l)
	h.l = l

	return h
}

func (h *handler) endpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		function, srv, ctxV, reqV, err := h.readRequest(ctx, service, req)
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

func (h *handler) readRequest(ctx context.Context, service Service, req any) (function, srv, ctxV, reqV reflect.Value, err error) {
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

func (h *handler) ServerWater(ctx context.Context, req any) (resp any, err error) {
	if len(h.finalizer) > 0 {
		defer func() {
			for _, fn := range h.finalizer {
				fn(ctx, err)
			}
		}()
	}

	resp, err = h.e(ctx, req)
	if err != nil {
		h.l.Error(err.Error())
		return nil, err
	}

	return resp, nil
}

func (h *handler) GetLogger() *slog.Logger {
	return h.l
}
