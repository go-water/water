package water

import "go.uber.org/zap"

type Service interface {
	Endpoint() Endpoint
	Name(srv Service) string
	SetLogger(l *zap.Logger)
}
