package water

import "go.uber.org/zap"

type Service interface {
	Endpoint() Endpoint
	Name() string
	SetLogger(l *zap.Logger)
	GetRequest() any
}
