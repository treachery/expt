package expt

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_parseSpec(t *testing.T) {
	fmt.Println(parseSpec("1-10"))
	fmt.Println(parseSpec("0-0"))
	fmt.Println(parseSpec("99-100"))
	fmt.Println(parseSpec("dff"))
	fmt.Println(parseSpec("10-8"))
}

func Test_layer(t *testing.T) {
	/*
		|-------------------|-------------------|
		|     expt2:50%   	|               	|
		|---------|---------|               	|
		|expt3:25%|expt4:25%|               	|
		|---------|---------|    expt1:50%      |
		|expt6:25%|         |       			|
		|---------|expt5:25%|              		|
		|expt7:25%|         |              		|
		|---------|---------|-------------------|
	*/
	root := &Layer{
		Id: "root",
		ModSpecToExpt: map[string]uint32{
			"50-99": 1,
		},
		ModSpecToChildLayers: map[string][]*Layer{
			"0-49": {
				&Layer{
					Id: "child-layer1",
					ModSpecToExpt: map[string]uint32{
						"0-99": 2,
					},
					ModSpecToChildLayers: map[string][]*Layer{},
				},
				&Layer{
					Id: "child-layer2",
					ModSpecToExpt: map[string]uint32{
						"0-49":  3,
						"50-99": 4,
					},
					ModSpecToChildLayers: map[string][]*Layer{},
				},
				&Layer{
					Id: "child-layer3",
					ModSpecToExpt: map[string]uint32{
						"50-99": 5,
					},
					ModSpecToChildLayers: map[string][]*Layer{
						"0-49": {
							&Layer{
								Id: "child-layer3-child-layer1",
								ModSpecToExpt: map[string]uint32{
									"0-99": 6,
								},
								ModSpecToChildLayers: map[string][]*Layer{},
							},
							&Layer{
								Id: "child-layer3-child-layer2",
								ModSpecToExpt: map[string]uint32{
									"0-99": 7,
								},
								ModSpecToChildLayers: map[string][]*Layer{},
							},
						},
					},
				},
			},
		},
	}
	bs, err := json.Marshal(root)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
	if err := root.parseModId(); err != nil {
		panic(err)
	}

	mexptcount := make(map[uint32]int)
	for i := 1; i <= 10000; i++ {
		ids, err := root.GetExptByHashId(fmt.Sprint(i))
		if err != nil {
			panic(err)
		}
		for _, id := range ids {
			mexptcount[id]++
		}
	}
	fmt.Println(mexptcount)
}
