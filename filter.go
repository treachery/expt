package expt

import (
	"context"
)

// 分层过滤
const (
	LAYER_FILTER_WHITE = uint8(1)
	LAYER_FILTER_BLACK = uint8(2)
)

type LayerFilter interface {
	Apply(key string) []uint32
}

type WhiteLayerFilter struct {
	whiteConfig map[string][]uint32
}

func (f *WhiteLayerFilter) Apply(key string) []uint32 {
	return f.whiteConfig[key]
}

// 实验过滤
const (
	EXPT_FILTER_GRAY  = uint8(3)
	EXPT_FILTER_GROUP = uint8(4)
	EXPT_FILTER_RULE  = uint8(5)
)

type ExptFilter interface {
	Filter(ctx context.Context, key string) (bool, string)
}

type GrayExptFilter struct {
}

type GroupExptFilter struct {
}

type RuleExptFilter struct {
}
