package water

import (
	"github.com/go-water/water/logger"
	"log/slog"
)

var Logger *slog.Logger

func init() {
	Logger = logger.NewLogger()
}
