package water

import (
	"context"
	"google.golang.org/grpc/metadata"
)

// ServerRequestFunc may take information from the received metadata header and
// use it to place items in the request scoped context. ServerRequestFuncs are
// executed prior to invoking the endpoint.
type ServerRequestFunc func(context.Context, metadata.MD) context.Context

// ServerResponseFunc may take information from a request context and use it to
// manipulate the gRPC response metadata headers and trailers. ResponseFuncs are
// only executed in servers, after invoking the endpoint but prior to writing a
// response.
type ServerResponseFunc func(ctx context.Context, header *metadata.MD, trailer *metadata.MD) context.Context
