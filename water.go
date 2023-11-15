package water

import (
	"github.com/go-water/water/endpoint"
	"go.uber.org/zap"
)

type Service interface {
	Endpoint() endpoint.Endpoint
	Name(srv Service) string
	SetLogger(l *zap.Logger)
}
