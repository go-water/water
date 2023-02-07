package water

import "go.uber.org/zap"

type Service interface {
	Endpoint() Endpoint
	Name() string
	SetLogger(l *zap.Logger)
}

func NewHandler(svc Service, options ...ServerOption) *Server {
	return NewServer(svc, options...)
}
