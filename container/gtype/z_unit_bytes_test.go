// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gtype_test

import (
	"github.com/basicfu/gf/internal/json"
	"github.com/basicfu/gf/util/gconv"
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_Bytes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		i := NewBytes([]byte("abc"))
		iClone := i.Clone()
		t.AssertEQ(iClone.Set([]byte("123")), []byte("abc"))
		t.AssertEQ(iClone.Val(), []byte("123"))

		//空参测试
		i1 := NewBytes()
		t.AssertEQ(i1.Val(), nil)
	})
}

func Test_Bytes_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		b := []byte("i love gf")
		i := NewBytes(b)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := NewBytes()
		err := json.Unmarshal(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), b)
	})
}

func Test_Bytes_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *Bytes
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
