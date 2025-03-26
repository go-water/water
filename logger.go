package water

import (
	"github.com/go-water/water/logger"
	"log/slog"
)

var (
	log = logger.NewLogger(logger.Level, logger.AddSource)
)

func NewLogger() *slog.Logger {
	return log
}

func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	log.Error(msg, args...)
}
