package water

import "go.uber.org/zap"

type Service interface {
	Endpoint() Endpoint
	GetServiceName(srv Service) string
	SetLogger(l *zap.Logger)
	GetRequest() any
}
