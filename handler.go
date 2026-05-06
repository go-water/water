package water

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/go-water/water/circuitbreaker"
	"github.com/go-water/water/endpoint"
	"github.com/go-water/water/logger"
	"github.com/sony/gobreaker"
)

type HandlerFunc func(*Context)

type RouterHandler struct {
	wt *Water
	h  HandlerFunc
}

func (r *RouterHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := r.wt.pool.Get().(*Context)
	ctx.Writer = w
	ctx.Request = req
	ctx.wt = r.wt
	ctx.reset()

	r.h(ctx)
	r.wt.pool.Put(ctx)
}

type handler struct {
	e         endpoint.Endpoint
	filter    Filter
	finalizer []FinalizerFunc
	l         *slog.Logger

	// 全局限流
	dl endpoint.Middleware
	el endpoint.Middleware
	// IP限流
	dil endpoint.Middleware
	eil endpoint.Middleware
	// 用户限流
	dul endpoint.Middleware
	eul endpoint.Middleware

	// 熔断
	breaker *gobreaker.CircuitBreaker
}

func NewHandler(srv Service, options ...ServerOption) Handler {
	h := new(handler)
	for _, option := range options {
		option(h)
	}

	h.e = h.endpoint(srv)
	if h.dl != nil {
		h.e = h.dl(h.e)
	}
	if h.el != nil {
		h.e = h.el(h.e)
	}

	if h.dil != nil {
		h.e = h.dil(h.e)
	}
	if h.eil != nil {
		h.e = h.eil(h.e)
	}

	if h.dul != nil {
		h.e = h.dul(h.e)
	}
	if h.eul != nil {
		h.e = h.eul(h.e)
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
		}

		err, ok := returnValue.(error)
		if ok {
			return nil, err
		}

		return nil, errors.New("method Handle return argument not include error type")
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
			err = errors.New("method Handle does not include three parameters")
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

	if h.filter != nil {
		err = h.filter(ctx)
		if err != nil {
			return nil, err
		}
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
