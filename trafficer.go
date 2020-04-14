package expt

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
)

// 分流逻辑
const (
	TRAFFIC_HASH       = uint8(1)
	TRAFFIC_TIMESLICE  = uint8(2)
	TRAFFIC_COUNTSLICE = uint8(3)
)

type Trafficer interface {
	Traffic(ctx context.Context, key string) (uint32, string, error)
}

type HashTrafficer struct {
	prefix string
	mods   map[uint32]uint32
}

func MustNewHashTrafficer(prefix string, modspec map[string]uint32) *HashTrafficer {
	t := &HashTrafficer{
		prefix: prefix,
		mods:   make(map[uint32]uint32),
	}
	for spec, vid := range modspec {
		low, high, err := parseModSpec(spec)
		if err != nil {
			panic(err)
		}
		for mod := low; mod <= high; mod++ {
			t.mods[uint32(mod)] = vid
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
		if _, err := hash.Write([]byte(fmt.Sprintf("%s%s", t.prefix, key))); err != nil {
			return 0, "hash error", err
		}
		mod := hash.Sum32() % 100
		return t.mods[mod], "HashTrafficer", nil
	}
}
