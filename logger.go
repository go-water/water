package water

import "github.com/go-water/water/logger"

var (
	l = logger.NewLogger(logger.Level, logger.AddSource)
)

func Info(msg string, args ...any) {
	l.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	l.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	l.Error(msg, args...)
}
