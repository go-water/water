package water

import (
	"context"
	"log/slog"
)

type ErrorHandler interface {
	Handle(ctx context.Context, err error)
}

type LogErrorHandler struct {
	l *slog.Logger
}

func NewLogErrorHandler(l *slog.Logger, n string) *LogErrorHandler {
	return &LogErrorHandler{
		l: l.WithGroup(n),
	}
}

func (h *LogErrorHandler) Handle(ctx context.Context, err error) {
	h.l.Error(err.Error())
}

func (h *LogErrorHandler) GetLogger() *slog.Logger {
	return h.l
}

type ErrorHandlerFunc func(ctx context.Context, err error)

func (f ErrorHandlerFunc) Handle(ctx context.Context, err error) {
	f(ctx, err)
}
