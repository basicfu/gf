// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gjson_test

import (
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_NewWithTag(t *testing.T) {
	type User struct {
		Age  int    `xml:"age-xml"  json:"age-json"`
		Name string `xml:"name-xml" json:"name-json"`
		Addr string `xml:"addr-xml" json:"addr-json"`
	}
	data := User{
		Age:  18,
		Name: "john",
		Addr: "chengdu",
	}
	// JSON
	gtest.C(t, func(t *gtest.T) {
		j := New(data)
		t.AssertNE(j, nil)
		t.Assert(j.Get("age-xml"), nil)
		t.Assert(j.Get("age-json"), data.Age)
		t.Assert(j.Get("name-xml"), nil)
		t.Assert(j.Get("name-json"), data.Name)
		t.Assert(j.Get("addr-xml"), nil)
		t.Assert(j.Get("addr-json"), data.Addr)
	})
	// XML
	gtest.C(t, func(t *gtest.T) {
		j := NewWithTag(data, "xml")
		t.AssertNE(j, nil)
		t.Assert(j.Get("age-xml"), data.Age)
		t.Assert(j.Get("age-json"), nil)
		t.Assert(j.Get("name-xml"), data.Name)
		t.Assert(j.Get("name-json"), nil)
		t.Assert(j.Get("addr-xml"), data.Addr)
		t.Assert(j.Get("addr-json"), nil)
	})
}

func Test_New_CustomStruct(t *testing.T) {
	type Base struct {
		Id int
	}
	type User struct {
		Base
		Name string
	}
	user := new(User)
	user.Id = 1
	user.Name = "john"

	gtest.C(t, func(t *gtest.T) {
		j := New(user)
		t.AssertNE(j, nil)

		s, err := j.ToJsonString()
		t.Assert(err, nil)
		t.Assert(s == `{"Uid":1,"Name":"john"}` || s == `{"Name":"john","Uid":1}`, true)
	})
}

func Test_New_HierarchicalStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Me struct {
			Name     string `json:"name"`
			Score    int    `json:"score"`
			Children []Me   `json:"children"`
		}
		me := Me{
			Name:  "john",
			Score: 100,
			Children: []Me{
				{
					Name:  "Bean",
					Score: 99,
				},
				{
					Name:  "Sam",
					Score: 98,
				},
			},
		}
		j := New(me)
		t.Assert(j.Remove("children.0.score"), nil)
		t.Assert(j.Remove("children.1.score"), nil)
		t.Assert(j.MustToJsonString(), `{"children":[{"children":null,"name":"Bean"},{"children":null,"name":"Sam"}],"name":"john","score":100}`)
	})
}

func Test_NewWithOption(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := []byte("[9223372036854775807, 9223372036854775806]")
		array := New(data).Array()
		t.Assert(array, []uint64{9223372036854776000, 9223372036854776000})
	})
	gtest.C(t, func(t *gtest.T) {
		data := []byte("[9223372036854775807, 9223372036854775806]")
		array := NewWithOption(data, Option{StrNumber: true}).Array()
		t.Assert(array, []uint64{9223372036854775807, 9223372036854775806})
	})
}
