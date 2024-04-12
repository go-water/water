package water

import "github.com/go-water/water/logger"

var (
	log = logger.NewLogger(logger.Level, logger.AddSource)
)

func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	log.Error(msg, args...)
}
