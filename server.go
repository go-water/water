package water

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-water/water/circuitbreaker"
	"github.com/go-water/water/endpoint"
	"github.com/go-water/water/ratelimit"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"reflect"
)

type Handler interface {
	ServerWater(ctx context.Context, req any) (any, error)
	GetLogger() *zap.Logger
}

type Server struct {
	e            endpoint.Endpoint
	finalizer    []ServerFinalizerFunc
	errorHandler ErrorHandler
	l            *zap.Logger
	limit        *rate.Limiter
	breaker      *gobreaker.CircuitBreaker
}

func NewHandler(srv Service, options ...ServerOption) Handler {
	s := new(Server)
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

	handler := NewLogErrorHandler(log, srv.Name(srv))
	srv.SetLogger(handler.l)
	s.l = handler.l
	s.errorHandler = handler

	return s
}

func (s *Server) endpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		function, srv, ctxV, reqV, err := s.readRequest(ctx, service, req)
		if err != nil {
			return nil, err
		}

		returnValues := function.Call([]reflect.Value{srv, ctxV, reflect.ValueOf(reqV.Interface())})
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

func (s *Server) readRequest(ctx context.Context, service Service, req any) (function, srv, ctxV, reqV reflect.Value, err error) {
	typ := reflect.TypeOf(service)
	srv = reflect.ValueOf(service)

	method, ok := typ.MethodByName("Handle")
	if ok {
		mType := method.Type
		num := mType.NumIn()
		if num == 3 {
			argType := mType.In(2)
			ctxV = reflect.ValueOf(ctx)
			reqV = reflect.New(argType.Elem())
			function = method.Func
			err = s.decodeRequest(req, reqV.Interface())
		} else {
			err = errors.New("method Handle does not include two parameters")
		}
	} else {
		err = errors.New("method Handle not implemented")
	}

	return
}

func (s *Server) decodeRequest(r, v any) (err error) {
	buf, err := json.Marshal(r)
	if err == nil {
		err = json.Unmarshal(buf, v)
	}

	return err
}

func (s *Server) ServerWater(ctx context.Context, req any) (resp any, err error) {
	if len(s.finalizer) > 0 {
		defer func() {
			for _, fn := range s.finalizer {
				fn(ctx, err)
			}
		}()
	}

	resp, err = s.e(ctx, req)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		return nil, err
	}

	return resp, nil
}

func (s *Server) GetLogger() *zap.Logger {
	return s.l
}
