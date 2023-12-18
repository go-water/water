package water

import (
	"go.uber.org/zap"
)

type Service interface {
	Name(srv Service) string
	SetLogger(l *zap.Logger)
}
