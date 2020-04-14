package expt

import (
	"encoding/json"
	"fmt"
	"hash/fnv"

	"github.com/pkg/errors"
)

var ErrSpec = errors.New("invalid spec")

//ModSpec 为闭区间如 [0-50][51-99]
type Layer struct {
	Id                   string              `json:"id"`
	ModSpecToExpt        map[string]uint32   `json:"modSpecToExpt"`
	ModSpecToChildLayers map[string][]*Layer `json:"modSpecToChildLayers"`
	modIdToExpt          map[uint32]uint32
	modIdToChildLayers   map[uint32][]*Layer
}

func NewLayerFromSPEC(spec string) (*Layer, error) {
	l := &Layer{}
	if err := json.Unmarshal([]byte(spec), l); err != nil {
		return nil, errors.Wrapf(err, "spec(%s)", spec)
	}
	if err := l.parseModId(); err != nil {
		return nil, err
	}
	return l, nil
}

func (root *Layer) parseModId() error {
	root.modIdToExpt = make(map[uint32]uint32)
	root.modIdToChildLayers = make(map[uint32][]*Layer)
	mmod := make(map[int]struct{})
	checkmod := func(mod int) error {
		if _, ok := mmod[mod]; ok {
			return errors.Errorf("mod(%d) confict", mod)
		}
		mmod[mod] = struct{}{}
		return nil
	}
	for spec, exptid := range root.ModSpecToExpt {
		low, high, err := parseModSpec(spec)
		if err != nil {
			return errors.Wrapf(err, "%s", spec)
		}
		for mod := low; mod <= high; mod++ {
			if err := checkmod(mod); err != nil {
				return err
			}
			root.modIdToExpt[uint32(mod)] = exptid
		}
	}
	for spec, childlayers := range root.ModSpecToChildLayers {
		for _, child := range childlayers {
			if err := child.parseModId(); err != nil {
				return err
			}
		}
		low, high, err := parseModSpec(spec)
		if err != nil {
			return errors.Wrapf(err, "%s", spec)
		}
		for mod := low; mod <= high; mod++ {
			if err := checkmod(mod); err != nil {
				return err
			}
			root.modIdToChildLayers[uint32(mod)] = childlayers
		}
	}
	return nil
}

func (root *Layer) GetExptByHashId(hashid string) (exptids []uint32, err error) {
	hash := fnv.New32()
	if _, err = hash.Write([]byte(fmt.Sprintf("%s_%s", root.Id, hashid))); err != nil {
		return exptids, errors.WithStack(err)
	}
	mod := hash.Sum32() % 100
	if exptid, ok := root.modIdToExpt[mod]; ok {
		return []uint32{exptid}, nil
	}
	if childs, ok := root.modIdToChildLayers[mod]; ok {
		for _, child := range childs {
			ids, err := child.GetExptByHashId(hashid)
			if err != nil {
				return exptids, errors.WithStack(err)
			}
			exptids = append(exptids, ids...)
		}
	}
	return exptids, nil
}
