package water

import (
	"context"
	"google.golang.org/grpc/grpclog"
)

// ErrorHandler receives a transport error to be processed for diagnostic purposes.
// Usually this means logging the error.
type ErrorHandler interface {
	Handle(ctx context.Context, err error)
}

// LogErrorHandler is a transport error handler implementation which logs an error.
type LogErrorHandler struct {
	logger grpclog.DepthLoggerV2
}

func NewLogErrorHandler(com string) *LogErrorHandler {
	return &LogErrorHandler{
		logger: grpclog.Component(com),
	}
}

func (h *LogErrorHandler) Handle(ctx context.Context, err error) {
	h.logger.Errorf("err: %s", err.Error())
}

// The ErrorHandlerFunc type is an adapter to allow the use of
// ordinary function as ErrorHandler. If f is a function
// with the appropriate signature, ErrorHandlerFunc(f) is a
// ErrorHandler that calls f.
type ErrorHandlerFunc func(ctx context.Context, err error)

// Handle calls f(ctx, err).
func (f ErrorHandlerFunc) Handle(ctx context.Context, err error) {
	f(ctx, err)
}
