package water

import (
	"context"
	"errors"
	"github.com/go-water/water/circuitbreaker"
	"github.com/go-water/water/endpoint"
	"github.com/go-water/water/ratelimit"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"reflect"
)

type adapter struct {
	e            endpoint.Endpoint
	finalizer    []ServerFinalizerFunc
	errorHandler ErrorHandler
	l            *zap.Logger
	limit        *rate.Limiter
	breaker      *gobreaker.CircuitBreaker
}

func NewAccessor(srv Service, options ...ServerOption) Accessor {
	s := new(adapter)
	for _, option := range options {
		option(s)
	}

	s.e = s.endpoint(srv)
	if s.limit != nil {
		s.e = ratelimit.NewErrorLimiter(s.limit)(s.e)
	}
	if s.breaker != nil {
		s.e = circuitbreaker.GoBreaker(s.breaker)(s.e)
	}

	handler := NewLogErrorHandler(Logger, srv.Name(srv))
	srv.SetLogger(handler.GetLogger())
	s.l = handler.GetLogger()
	s.errorHandler = handler

	return s
}

func (ad *adapter) endpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		function, srv, ctxV, reqV, err := ad.readRequest(ctx, service, req)
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

func (ad *adapter) readRequest(ctx context.Context, service Service, req any) (function, srv, ctxV, reqV reflect.Value, err error) {
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

func (ad *adapter) ServerWater(ctx context.Context, req any) (resp any, err error) {
	if len(ad.finalizer) > 0 {
		defer func() {
			for _, fn := range ad.finalizer {
				fn(ctx, err)
			}
		}()
	}

	resp, err = ad.e(ctx, req)
	if err != nil {
		ad.errorHandler.Handle(ctx, err)
		return nil, err
	}

	return resp, nil
}

func (ad *adapter) GetLogger() *zap.Logger {
	return ad.l
}
