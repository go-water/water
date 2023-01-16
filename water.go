package water

type Service interface {
	Endpoint() Endpoint
}

func NewHandler(svc Service, options ...ServerOption) *Server {
	return NewServer(
		svc.Endpoint(),
		options...,
	)
}
