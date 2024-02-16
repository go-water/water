package logger

import (
	"log/slog"
	"os"
)

func NewLogger(level slog.Leveler, addSource bool) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level, AddSource: addSource, ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		v := a.Value
		if v.Kind() == slog.KindTime {
			return slog.String(a.Key, v.Time().Format("2006-01-02 15:04:05"))
		}

		return a
	}}))
}
