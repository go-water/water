package water

import (
	"context"
	"go.uber.org/zap"
)

type ErrorHandler interface {
	Handle(ctx context.Context, err error)
}

type LogErrorHandler struct {
	n      string
	logger *zap.Logger
}

func NewLogErrorHandler(l *zap.Logger, n string) *LogErrorHandler {
	return &LogErrorHandler{
		logger: l,
		n:      n,
	}
}

func (h *LogErrorHandler) Handle(ctx context.Context, err error) {
	h.logger.Error(h.n, zap.Error(err))
}

type ErrorHandlerFunc func(ctx context.Context, err error)

func (f ErrorHandlerFunc) Handle(ctx context.Context, err error) {
	f(ctx, err)
}
