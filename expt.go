package expt

import (
	"context"
	"encoding/json"
)

type Expt struct {
	Id         uint32   `json:"id"`
	VersionIds []uint32 `json:"version_ids"`

	Selectors []Selector   `json:"selectors"`
	Filters   []ExptFilter `json:"filters"`
	Trafficer Trafficer    `json:"trafficer"`
}

func NewExpt(id uint32, vids []uint32, trafficer Trafficer, opts ...func(e *Expt)) *Expt {
	e := &Expt{
		Id:         id,
		VersionIds: vids,
		Trafficer:  trafficer,
	}
	for _, o := range opts {
		o(e)
	}
	return e
}

func (e *Expt) UnmarshalJSON(bs []byte) error {
	decoder := new(struct {
		Id         uint32   `json:"id"`
		VersionIds []uint32 `json:"version_ids"`

		Selectors []struct {
			SelectType string `json:"SelectType"`
			Obj        json.RawMessage
		} `json:"selectors"`
		Filters []struct {
			FilterType string `json:"FilterType"`
			Obj        json.RawMessage
		} `json:"filters"`
		Trafficer struct {
			TrafficType string          `json:"TrafficType"`
			Obj         json.RawMessage `json:"object"`
		} `json:"trafficer"`
	})
	err := json.Unmarshal(bs, decoder)
	if err == nil {
		e.Id = decoder.Id
		e.VersionIds = decoder.VersionIds
		for _, selector := range decoder.Selectors {
			switch selector.SelectType {
			//TODO
			}
		}
		for _, filter := range decoder.Filters {
			switch filter.FilterType {
			//TODO
			}
		}
		switch decoder.Trafficer.TrafficType {
		case TRAFFIC_HASH:
			trafficer := &HashTrafficer{
				TrafficType: TRAFFIC_HASH,
			}
			if err := json.Unmarshal(decoder.Trafficer.Obj, &trafficer.Object); err != nil {
				return err
			}
			e.Trafficer = trafficer
		}
		return nil
	}
	return err
}

func (e *Expt) Run(ctx context.Context, key string) (versionid uint32, msg string, err error) {
	for _, selector := range e.Selectors {
		versionid, msg, err = selector.Select(ctx, key)
		if err != nil {
			return
		}
		if versionid > 0 {
			return
		}
	}
	for _, filter := range e.Filters {
		var filterd bool
		if filterd, msg = filter.Filter(ctx, key); filterd {
			return
		}
	}
	return e.Trafficer.Traffic(ctx, key)
}
