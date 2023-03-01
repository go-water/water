package water

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type ServerBase struct {
	l *zap.Logger
}

func (s *ServerBase) Name(srv any) string {
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
