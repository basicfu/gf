// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gvar_test

import (
	"github.com/basicfu/gf/frame/g"
	"github.com/basicfu/gf/test/gtest"
	"github.com/basicfu/gf/util/gconv"
	"testing"
)

func TestVar_Struct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type StTest struct {
			Test int
		}

		Kv := make(map[string]int, 1)
		Kv["Test"] = 100

		testObj := &StTest{}

		objOne := New(Kv, true)

		objOne.Struct(testObj)

		t.Assert(testObj.Test, Kv["Test"])
	})
	gtest.C(t, func(t *gtest.T) {
		type StTest struct {
			Test int8
		}
		o := &StTest{}
		v := New(g.Slice{"Test", "-25"})
		v.Struct(o)
		t.Assert(o.Test, -25)
	})
}

func TestVar_Var_Attribute_Struct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid  int
			Name string
		}
		user := new(User)
		err := gconv.Struct(
			g.Map{
				"uid":  New(1),
				"name": New("john"),
			}, user)
		t.Assert(err, nil)
		t.Assert(user.Uid, 1)
		t.Assert(user.Name, "john")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid  int
			Name string
		}
		var user *User
		err := gconv.Struct(
			g.Map{
				"uid":  New(1),
				"name": New("john"),
			}, &user)
		t.Assert(err, nil)
		t.Assert(user.Uid, 1)
		t.Assert(user.Name, "john")
	})
}
