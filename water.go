package water

type Service interface {
	Endpoint() Endpoint
	//DecodeRequest(ctx context.Context, req interface{}) (interface{}, error)
	//EncodeResponse(ctx context.Context, resp interface{}) (interface{}, error)
	//ServerBefore() ServerOption
	//ServerAfter() ServerOption
	//GetLogger() grpclog.DepthLoggerV2
	//Name() string
}

func NewHandler(svc Service) *Server {
	return NewServer(
		svc.Endpoint(),
	)
}
