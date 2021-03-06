// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gconv_test

import (
	"testing"

	"github.com/basicfu/gf/g"
	"github.com/basicfu/gf/test/gtest"
)

func Test_MapToMap1(t *testing.T) {
	// map[int]int -> map[string]string
	// empty original map.
	gtest.C(t, func(t *gtest.T) {
		m1 := g.MapIntInt{}
		m2 := g.MapStrStr{}
		t.Assert(MapToMap(m1, &m2), nil)
		t.Assert(len(m1), len(m2))
	})
	// map[int]int -> map[string]string
	gtest.C(t, func(t *gtest.T) {
		m1 := g.MapIntInt{
			1: 100,
			2: 200,
		}
		m2 := g.MapStrStr{}
		t.Assert(MapToMap(m1, &m2), nil)
		t.Assert(m2["1"], m1[1])
		t.Assert(m2["2"], m1[2])
	})
	// map[string]interface{} -> map[string]string
	gtest.C(t, func(t *gtest.T) {
		m1 := g.Map{
			"k1": "v1",
			"k2": "v2",
		}
		m2 := g.MapStrStr{}
		t.Assert(MapToMap(m1, &m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
	// map[string]string -> map[string]interface{}
	gtest.C(t, func(t *gtest.T) {
		m1 := g.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		m2 := g.Map{}
		t.Assert(MapToMap(m1, &m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
	// map[string]interface{} -> map[interface{}]interface{}
	gtest.C(t, func(t *gtest.T) {
		m1 := g.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		m2 := g.MapAnyAny{}
		t.Assert(MapToMap(m1, &m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
}

func Test_MapToMap2(t *testing.T) {
	type User struct {
		Id   int
		Name string
	}
	params := g.Map{
		"key": g.Map{
			"id":   1,
			"name": "john",
		},
	}
	gtest.C(t, func(t *gtest.T) {
		m := make(map[string]User)
		err := MapToMap(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
	gtest.C(t, func(t *gtest.T) {
		m := (map[string]User)(nil)
		err := MapToMap(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
	gtest.C(t, func(t *gtest.T) {
		m := make(map[string]*User)
		err := MapToMap(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
	gtest.C(t, func(t *gtest.T) {
		m := (map[string]*User)(nil)
		err := MapToMap(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
}

func Test_MapToMapDeep(t *testing.T) {
	type Ids struct {
		Id  int
		Uid int
	}
	type Base struct {
		Ids
		Time string
	}
	type User struct {
		Base
		Name string
	}
	params := g.Map{
		"key": g.Map{
			"id":   1,
			"name": "john",
		},
	}
	gtest.C(t, func(t *gtest.T) {
		m := (map[string]*User)(nil)
		err := MapToMap(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
}

func Test_MapToMaps1(t *testing.T) {
	type User struct {
		Id   int
		Name int
	}
	params := g.Map{
		"key1": g.Slice{
			g.Map{"id": 1, "name": "john"},
			g.Map{"id": 2, "name": "smith"},
		},
		"key2": g.Slice{
			g.Map{"id": 3, "name": "green"},
			g.Map{"id": 4, "name": "jim"},
		},
	}
	gtest.C(t, func(t *gtest.T) {
		m := make(map[string][]User)
		err := MapToMaps(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 2)
		t.Assert(m["key1"][0].Id, 1)
		t.Assert(m["key1"][1].Id, 2)
		t.Assert(m["key2"][0].Id, 3)
		t.Assert(m["key2"][1].Id, 4)
	})
	gtest.C(t, func(t *gtest.T) {
		m := (map[string][]User)(nil)
		err := MapToMaps(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 2)
		t.Assert(m["key1"][0].Id, 1)
		t.Assert(m["key1"][1].Id, 2)
		t.Assert(m["key2"][0].Id, 3)
		t.Assert(m["key2"][1].Id, 4)
	})
	gtest.C(t, func(t *gtest.T) {
		m := make(map[string][]*User)
		err := MapToMaps(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 2)
		t.Assert(m["key1"][0].Id, 1)
		t.Assert(m["key1"][1].Id, 2)
		t.Assert(m["key2"][0].Id, 3)
		t.Assert(m["key2"][1].Id, 4)
	})
	gtest.C(t, func(t *gtest.T) {
		m := (map[string][]*User)(nil)
		err := MapToMaps(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 2)
		t.Assert(m["key1"][0].Id, 1)
		t.Assert(m["key1"][1].Id, 2)
		t.Assert(m["key2"][0].Id, 3)
		t.Assert(m["key2"][1].Id, 4)
	})
}

func Test_MapToMaps2(t *testing.T) {
	type User struct {
		Id   int
		Name int
	}
	params := g.MapIntAny{
		100: g.Slice{
			g.Map{"id": 1, "name": "john"},
			g.Map{"id": 2, "name": "smith"},
		},
		200: g.Slice{
			g.Map{"id": 3, "name": "green"},
			g.Map{"id": 4, "name": "jim"},
		},
	}
	gtest.C(t, func(t *gtest.T) {
		m := make(map[int][]User)
		err := MapToMaps(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 2)
		t.Assert(m[100][0].Id, 1)
		t.Assert(m[100][1].Id, 2)
		t.Assert(m[200][0].Id, 3)
		t.Assert(m[200][1].Id, 4)
	})
	gtest.C(t, func(t *gtest.T) {
		m := make(map[int][]*User)
		err := MapToMaps(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 2)
		t.Assert(m[100][0].Id, 1)
		t.Assert(m[100][1].Id, 2)
		t.Assert(m[200][0].Id, 3)
		t.Assert(m[200][1].Id, 4)
	})
	gtest.C(t, func(t *gtest.T) {
		m := make(map[string][]*User)
		err := MapToMaps(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 2)
		t.Assert(m["100"][0].Id, 1)
		t.Assert(m["100"][1].Id, 2)
		t.Assert(m["200"][0].Id, 3)
		t.Assert(m["200"][1].Id, 4)
	})
}

func Test_MapToMaps3(t *testing.T) {
	type Ids struct {
		Id  int
		Uid int
	}
	type Base struct {
		Ids
		Time string
	}
	type User struct {
		Base
		Name string
	}
	params := g.MapIntAny{
		100: g.Slice{
			g.Map{"id": 1, "name": "john"},
			g.Map{"id": 2, "name": "smith"},
		},
		200: g.Slice{
			g.Map{"id": 3, "name": "green"},
			g.Map{"id": 4, "name": "jim"},
		},
	}
	gtest.C(t, func(t *gtest.T) {
		m := make(map[string][]*User)
		err := MapToMaps(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 2)
		t.Assert(m["100"][0].Id, 1)
		t.Assert(m["100"][1].Id, 2)
		t.Assert(m["100"][0].Name, "john")
		t.Assert(m["100"][1].Name, "smith")
		t.Assert(m["200"][0].Id, 3)
		t.Assert(m["200"][1].Id, 4)
		t.Assert(m["200"][0].Name, "green")
		t.Assert(m["200"][1].Name, "jim")
	})
}

func Test_MapToMapsWithTag(t *testing.T) {
	type Ids struct {
		Id  int
		Uid int
	}
	type Base struct {
		Ids  `json:"ids"`
		Time string
	}
	type User struct {
		Base `json:"base"`
		Name string
	}
	params := g.MapIntAny{
		100: g.Slice{
			g.Map{"id": 1, "name": "john"},
			g.Map{"id": 2, "name": "smith"},
		},
		200: g.Slice{
			g.Map{"id": 3, "name": "green"},
			g.Map{"id": 4, "name": "jim"},
		},
	}
	gtest.C(t, func(t *gtest.T) {
		m := make(map[string][]*User)
		err := MapToMaps(params, &m)
		t.Assert(err, nil)
		t.Assert(len(m), 2)
		t.Assert(m["100"][0].Id, 1)
		t.Assert(m["100"][1].Id, 2)
		t.Assert(m["100"][0].Name, "john")
		t.Assert(m["100"][1].Name, "smith")
		t.Assert(m["200"][0].Id, 3)
		t.Assert(m["200"][1].Id, 4)
		t.Assert(m["200"][0].Name, "green")
		t.Assert(m["200"][1].Name, "jim")
	})
}
