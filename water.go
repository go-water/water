package water

type Service interface {
	Endpoint() Endpoint
	Name() string
}

func NewHandler(svc Service, options ...ServerOption) *Server {
	return NewServer(
		svc.Endpoint(),
		svc.Name(),
		options...,
	)
}
