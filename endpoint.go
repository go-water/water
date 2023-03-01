package water

import (
	"context"
)

type Endpoint func(ctx context.Context, req any) (any, error)
