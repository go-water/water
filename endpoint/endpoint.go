package endpoint

import (
	"context"
)

type Endpoint func(ctx context.Context, req any) (any, error)

type Middleware func(Endpoint) Endpoint
