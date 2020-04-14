package expt

import (
	"context"
)

type Selector interface {
	Select(ctx context.Context, key string) (uint32, string, error)
}

type WhiteSelector struct{}
type CacheSelector struct{}
