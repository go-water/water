package water

import (
	"github.com/go-water/water/logger"
	"log/slog"
)

var (
	Logger    *slog.Logger
	Level     = slog.LevelInfo
	AddSource bool
)

func init() {
	Logger = logger.NewLogger(Level, AddSource)
}
