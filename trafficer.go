package expt

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
)

// 分流逻辑
const (
	TRAFFIC_HASH       = "TRAFFIC_HASH"
	TRAFFIC_TIMESLICE  = "TRAFFIC_TIMESLICE"
	TRAFFIC_COUNTSLICE = "TRAFFIC_COUNTSLICE"
)

type Trafficer interface {
	Traffic(ctx context.Context, key string) (uint32, string, error)
}

// 定义为二层结构，是为了方便序列化，放到缓存中
type HashTrafficer struct {
	TrafficType string `json:"TrafficType"`
	Object      struct {
		Prefix string            `json:"prefix"`
		Mods   map[uint32]uint32 `json:"mods"`
	} `json:"object"`
}

func MustNewHashTrafficer(prefix string, modspec map[string]uint32) *HashTrafficer {
	t := &HashTrafficer{
		TrafficType: TRAFFIC_HASH,
	}
	t.Object.Prefix = prefix
	t.Object.Mods = make(map[uint32]uint32)
	for spec, vid := range modspec {
		low, high, err := parseModSpec(spec)
		if err != nil {
			panic(err)
		}
		for mod := low; mod <= high; mod++ {
			t.Object.Mods[uint32(mod)] = vid
		}
	}
	return t
}

func (t *HashTrafficer) Traffic(ctx context.Context, key string) (uint32, string, error) {
	select {
	case <-ctx.Done():
		return 0, "context deadline", errors.New("context deadline")
	default:
		hash := fnv.New32()
		if _, err := hash.Write([]byte(fmt.Sprintf("%s%s", t.Object.Prefix, key))); err != nil {
			return 0, "hash error", err
		}
		mod := hash.Sum32() % 100
		return t.Object.Mods[mod], t.TrafficType, nil
	}
}
