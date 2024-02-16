package water

import (
	"github.com/go-water/water/logger"
	"log/slog"
)

var Logger *slog.Logger

func InitLog() {
	Logger = logger.NewLogger()
}

func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}
