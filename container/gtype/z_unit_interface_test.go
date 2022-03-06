// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gtype_test

import (
	"github.com/basicfu/gf/json"
	"github.com/basicfu/gf/test/gtest"
	"github.com/basicfu/gf/util/gconv"
	"testing"
)

func Test_Interface(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t1 := Temp{Name: "gf", Age: 18}
		t2 := Temp{Name: "gf", Age: 19}
		i := New(t1)
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(t2), t1)
		t.AssertEQ(iClone.Val().(Temp), t2)

		//空参测试
		i1 := New()
		t.AssertEQ(i1.Val(), nil)
	})
}

func Test_Interface_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "i love gf"
		i := New(s)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := New()
		err := json.Unmarshal(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), s)
	})
}

func Test_Interface_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *Interface
	}
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), "123")
	})
}
