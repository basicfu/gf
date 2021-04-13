// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gtype_test

import (
	"github.com/basicfu/gf/internal/json"
	"github.com/basicfu/gf/test/gtest"
	"github.com/basicfu/gf/util/gconv"
	"math"
	"testing"
)

func Test_Float64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		i := NewFloat64(0)
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(0.1), float64(0))
		t.AssertEQ(iClone.Val(), float64(0.1))
		//空参测试
		i1 := NewFloat64()
		t.AssertEQ(i1.Val(), float64(0))
	})
}

func Test_Float64_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := math.MaxFloat64
		i := NewFloat64(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := NewFloat64()
		err := json.Unmarshal(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), v)
	})
}

func Test_Float64_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *Float64
	}
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123.456",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), "123.456")
	})
}
