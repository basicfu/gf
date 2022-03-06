// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gutil_test

import (
	"github.com/basicfu/gf/g"
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_Dump(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		Dump(map[int]int{
			100: 100,
		})
	})
}

func Test_Try(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := `gutil Try test`
		t.Assert(Try(func() {
			panic(s)
		}), s)
	})
}

func Test_TryCatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		TryCatch(func() {
			panic("gutil TryCatch test")
		})
	})

	gtest.C(t, func(t *gtest.T) {
		TryCatch(func() {
			panic("gutil TryCatch test")

		}, func(err error) {
			t.Assert(err, "gutil TryCatch test")
		})
	})
}

func Test_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(IsEmpty(1), false)
	})
}

func Test_Throw(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			t.Assert(recover(), "gutil Throw test")
		}()

		Throw("gutil Throw test")
	})
}

func Test_Keys(t *testing.T) {
	// map
	gtest.C(t, func(t *gtest.T) {
		keys := Keys(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)
	})
	// *map
	gtest.C(t, func(t *gtest.T) {
		keys := Keys(&map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)
	})
	// *struct
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		keys := Keys(new(T))
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
	// *struct nil
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := Keys(pointer)
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
	// **struct nil
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := Keys(&pointer)
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
}

func Test_Values(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		values := Keys(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", values)
		t.AssertIN("2", values)
	})

	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		keys := Values(T{
			A: "1",
			B: 2,
		})
		t.Assert(keys, g.Slice{"1", 2})
	})
}
