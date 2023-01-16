package water

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ErrorHandler interface {
	Handle(ctx *gin.Context, err error)
}

type LogErrorHandler struct {
	logger *zap.Logger
}

func NewLogErrorHandler(l *zap.Logger) *LogErrorHandler {
	return &LogErrorHandler{
		logger: l,
	}
}

func (h *LogErrorHandler) Handle(ctx *gin.Context, err error) {
	h.logger.Error(fmt.Sprintf("err: %s", err.Error()))
}

type ErrorHandlerFunc func(ctx *gin.Context, err error)

func (f ErrorHandlerFunc) Handle(ctx *gin.Context, err error) {
	f(ctx, err)
}
