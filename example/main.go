package main

import (
	"context"
	"fmt"

	"github.com/treachery/expt"
)

/*
分层模型如下:
为了方便，每个实验默认两个版本versionid=1,2
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

var expts map[uint32]*expt.Expt

func initExpts() {
	expts = make(map[uint32]*expt.Expt)
	expts[1] = expt.NewExpt(1, []uint32{1, 2},
		expt.WithTrafficer(expt.MustNewHashTrafficer("prefix_", map[string]uint32{
			"0-49":  1,
			"50-99": 2,
		})))
	expts[2] = expt.NewExpt(2, []uint32{1, 2},
		expt.WithTrafficer(expt.MustNewHashTrafficer("prefix_", map[string]uint32{
			"0-49":  1,
			"50-99": 2,
		})))
	expts[3] = expt.NewExpt(3, []uint32{1, 2},
		expt.WithTrafficer(expt.MustNewHashTrafficer("prefix_", map[string]uint32{
			"0-49":  1,
			"50-99": 2,
		})))
	expts[4] = expt.NewExpt(4, []uint32{1, 2},
		expt.WithTrafficer(expt.MustNewHashTrafficer("prefix_", map[string]uint32{
			"0-49":  1,
			"50-99": 2,
		})))
	expts[5] = expt.NewExpt(5, []uint32{1, 2},
		expt.WithTrafficer(expt.MustNewHashTrafficer("prefix_", map[string]uint32{
			"0-49":  1,
			"50-99": 2,
		})))
	expts[6] = expt.NewExpt(6, []uint32{1, 2},
		expt.WithTrafficer(expt.MustNewHashTrafficer("prefix_", map[string]uint32{
			"0-49":  1,
			"50-99": 2,
		})))
	expts[7] = expt.NewExpt(7, []uint32{1, 2},
		expt.WithTrafficer(expt.MustNewHashTrafficer("prefix_", map[string]uint32{
			"0-49":  1,
			"50-99": 2,
		})))
}

var spec = `{
		"id": "root",
		"modSpecToExpt": {
		  "50-99": 1
		},
		"modSpecToChildLayers": {
		  "0-49": [
			{
			  "id": "child-layer1",
			  "modSpecToExpt": {
				"0-99": 2
			  },
			  "modSpecToChildLayers": {}
			},
			{
			  "id": "child-layer2",
			  "modSpecToExpt": {
				"0-49": 3,
				"50-99": 4
			  },
			  "modSpecToChildLayers": {}
			},
			{
			  "id": "child-layer3",
			  "modSpecToExpt": {
				"50-99": 5
			  },
			  "modSpecToChildLayers": {
				"0-49": [
				  {
					"id": "child-layer3-child-layer1",
					"modSpecToExpt": {
					  "0-99": 6
					},
					"modSpecToChildLayers": {}
				  },
				  {
					"id": "child-layer3-child-layer2",
					"modSpecToExpt": {
					  "0-99": 7
					},
					"modSpecToChildLayers": {}
				  }
				]
			  }
			}
		  ]
		}
	  }`

func main() {
	initExpts()
	root, err := expt.ReadFromSPEC(spec)
	if err != nil {
		panic(err)
	}

	for i := 1; i < 100; i++ {
		key := fmt.Sprint(i)
		exptids, err := root.GetExptByHashId(key)
		if err != nil {
			panic(err)
		}
		for _, exptid := range exptids {
			vid, msg, err := expts[exptid].Run(context.Background(), key)
			fmt.Println(exptid, vid, msg, err)
		}
	}
}
