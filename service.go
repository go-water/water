package water

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

type Handler interface {
	ServerWater(ctx context.Context, req any) (any, error)
	GetLogger() *slog.Logger
}

type Service interface {
	Name(srv Service) string
	SetLogger(l *slog.Logger)
}

type ServerBase struct {
	l *slog.Logger
}

func (s *ServerBase) Name(srv Service) string {
	fullName := fmt.Sprintf("%T", srv)
	index := strings.LastIndex(fullName, ".")
	name := fullName[index+1:]
	return name
}

func (s *ServerBase) GetLogger() *slog.Logger {
	return s.l
}

func (s *ServerBase) SetLogger(l *slog.Logger) {
	s.l = l
}
