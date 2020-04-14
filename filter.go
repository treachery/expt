package expt

import (
	"context"
)

// 实验过滤
const (
	EXPT_FILTER_WHITE = uint8(1)
	EXPT_FILTER_BLACK = uint8(2)
)

type LayerFilter interface {
	Apply(key string) []uint32
}

type ExptFilterWhite struct {
	whiteConfig map[string][]uint32
}

func (f *ExptFilterWhite) Apply(key string) []uint32 {
	return f.whiteConfig[key]
}

// 分流过滤
const (
	TRAFFIC_FILTER_GRAY  = uint8(3)
	TRAFFIC_FILTER_GROUP = uint8(4)
	TRAFFIC_FILTER_RULE  = uint8(5)
)

type ExptFilter interface {
	Filter(ctx context.Context, key string) (bool, string)
}
