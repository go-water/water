package water

import (
	"github.com/go-water/water/logger"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

func init() {
	zLog := logger.Config{}
	grpclog.SetLoggerV2(zapgrpc.NewLogger(zLog.NewLogger()))
}
