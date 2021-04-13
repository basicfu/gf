// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gconv_test

import (
	"testing"

	"github.com/basicfu/gf/frame/g"
	"github.com/basicfu/gf/test/gtest"
)

func Test_Slice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value := 123.456
		t.AssertEQ(Bytes("123"), []byte("123"))
		t.AssertEQ(Strings(value), []string{"123.456"})
		t.AssertEQ(Ints(value), []int{123})
		t.AssertEQ(Floats(value), []float64{123.456})
		t.AssertEQ(Interfaces(value), []interface{}{123.456})
	})
}

func Test_Slice_Empty(t *testing.T) {
	// Int.
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(Ints(""), []int{})
		t.Assert(Ints(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(Int32s(""), []int32{})
		t.Assert(Int32s(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(Int64s(""), []int64{})
		t.Assert(Int64s(nil), nil)
	})
	// Uint.
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(Uints(""), []uint{})
		t.Assert(Uints(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(Uint32s(""), []uint32{})
		t.Assert(Uint32s(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(Uint64s(""), []uint64{})
		t.Assert(Uint64s(nil), nil)
	})
	// Float.
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(Floats(""), []float64{})
		t.Assert(Floats(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(Float32s(""), []float32{})
		t.Assert(Float32s(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(Float64s(""), []float64{})
		t.Assert(Float64s(nil), nil)
	})
}

func Test_Strings(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := []*g.Var{
			g.NewVar(1),
			g.NewVar(2),
			g.NewVar(3),
		}
		t.AssertEQ(Strings(array), []string{"1", "2", "3"})
	})
}

func Test_Slice_Interfaces(t *testing.T) {
	// map
	gtest.C(t, func(t *gtest.T) {
		array := Interfaces(g.Map{
			"id":   1,
			"name": "john",
		})
		t.Assert(len(array), 1)
		t.Assert(array[0].(g.Map)["id"], 1)
		t.Assert(array[0].(g.Map)["name"], "john")
	})
	// struct
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Id   int `json:"id"`
			Name string
		}
		array := Interfaces(&A{
			Id:   1,
			Name: "john",
		})
		t.Assert(len(array), 1)
		t.Assert(array[0].(*A).Id, 1)
		t.Assert(array[0].(*A).Name, "john")
	})
}

func Test_Slice_PrivateAttribute(t *testing.T) {
	type User struct {
		Id   int    `json:"id"`
		name string `json:"name"`
	}
	gtest.C(t, func(t *gtest.T) {
		user := &User{1, "john"}
		array := Interfaces(user)
		t.Assert(len(array), 1)
		t.Assert(array[0].(*User).Id, 1)
		t.Assert(array[0].(*User).name, "john")
	})
}

func Test_Slice_Structs(t *testing.T) {
	type Base struct {
		Age int
	}
	type User struct {
		Id   int
		Name string
		Base
	}

	gtest.C(t, func(t *gtest.T) {
		users := make([]User, 0)
		params := []g.Map{
			{"id": 1, "name": "john", "age": 18},
			{"id": 2, "name": "smith", "age": 20},
		}
		err := Structs(params, &users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, params[0]["id"])
		t.Assert(users[0].Name, params[0]["name"])
		t.Assert(users[0].Age, 18)

		t.Assert(users[1].Id, params[1]["id"])
		t.Assert(users[1].Name, params[1]["name"])
		t.Assert(users[1].Age, 20)
	})
}
