package water

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type Handler interface {
	ServerWater(ctx context.Context, req any) (any, error)
	GetLogger() *zap.Logger
}

type Service interface {
	Name(srv Service) string
	SetLogger(l *zap.Logger)
}

type ServerBase struct {
	l *zap.Logger
}

func (s *ServerBase) Name(srv Service) string {
	fullName := fmt.Sprintf("%T", srv)
	index := strings.LastIndex(fullName, ".")
	name := fullName[index+1:]
	return name
}

func (s *ServerBase) GetLogger() *zap.Logger {
	return s.l
}

func (s *ServerBase) SetLogger(l *zap.Logger) {
	s.l = l
}
