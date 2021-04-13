// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gconv_test

import (
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		f32 := float32(123.456)
		i64 := int64(1552578474888)
		t.AssertEQ(Int(f32), int(123))
		t.AssertEQ(Int8(f32), int8(123))
		t.AssertEQ(Int16(f32), int16(123))
		t.AssertEQ(Int32(f32), int32(123))
		t.AssertEQ(Int64(f32), int64(123))
		t.AssertEQ(Int64(f32), int64(123))
		t.AssertEQ(Uint(f32), uint(123))
		t.AssertEQ(Uint8(f32), uint8(123))
		t.AssertEQ(Uint16(f32), uint16(123))
		t.AssertEQ(Uint32(f32), uint32(123))
		t.AssertEQ(Uint64(f32), uint64(123))
		t.AssertEQ(Float32(f32), float32(123.456))
		t.AssertEQ(Float64(i64), float64(i64))
		t.AssertEQ(Bool(f32), true)
		t.AssertEQ(String(f32), "123.456")
		t.AssertEQ(String(i64), "1552578474888")
	})

	gtest.C(t, func(t *gtest.T) {
		s := "-0xFF"
		t.Assert(Int(s), int64(-0xFF))
	})
}

func Test_Duration(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		d := Duration("1s")
		t.Assert(d.String(), "1s")
		t.Assert(d.Nanoseconds(), 1000000000)
	})
}
