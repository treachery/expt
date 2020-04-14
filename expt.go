package expt

import (
	"context"
)

type Expt struct {
	Id         uint32
	VersionIds []uint32

	selectors []Selector
	filters   []ExptFilter
	trafficer Trafficer
}

func NewExpt(id uint32, vids []uint32, opts ...func(e *Expt)) *Expt {
	e := &Expt{
		Id:         id,
		VersionIds: vids,
	}
	for _, o := range opts {
		o(e)
	}
	return e
}

func (e *Expt) Run(ctx context.Context, key string) (versionid uint32, msg string, err error) {
	for _, selector := range e.selectors {
		versionid, msg, err = selector.Select(ctx, key)
		if err != nil {
			return
		}
		if versionid > 0 {
			return
		}
	}
	for _, filter := range e.filters {
		var filterd bool
		if filterd, msg = filter.Filter(ctx, key); filterd {
			return
		}
	}
	return e.trafficer.Traffic(ctx, key)
}
