package water

import (
	"context"
	"go.uber.org/zap"
)

type ErrorHandler interface {
	Handle(ctx context.Context, err error)
}

type LogErrorHandler struct {
	l *zap.Logger
}

func NewLogErrorHandler(l *zap.Logger, n string) *LogErrorHandler {
	return &LogErrorHandler{
		l: l.Named(n),
	}
}

func (h *LogErrorHandler) Handle(ctx context.Context, err error) {
	h.l.Error("Core", zap.Error(err))
}

type ErrorHandlerFunc func(ctx context.Context, err error)

func (f ErrorHandlerFunc) Handle(ctx context.Context, err error) {
	f(ctx, err)
}
