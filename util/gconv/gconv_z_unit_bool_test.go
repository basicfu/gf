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

type boolStruct struct {
}

func Test_Bool(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var i interface{} = nil
		t.AssertEQ(Bool(i), false)
		t.AssertEQ(Bool(false), false)
		t.AssertEQ(Bool(nil), false)
		t.AssertEQ(Bool(0), false)
		t.AssertEQ(Bool("0"), false)
		t.AssertEQ(Bool(""), false)
		t.AssertEQ(Bool("false"), false)
		t.AssertEQ(Bool("off"), false)
		t.AssertEQ(Bool([]byte{}), false)
		t.AssertEQ(Bool([]string{}), false)
		t.AssertEQ(Bool([]interface{}{}), false)
		t.AssertEQ(Bool([]map[int]int{}), false)

		t.AssertEQ(Bool("1"), true)
		t.AssertEQ(Bool("on"), true)
		t.AssertEQ(Bool(1), true)
		t.AssertEQ(Bool(123.456), true)
		t.AssertEQ(Bool(boolStruct{}), true)
		t.AssertEQ(Bool(&boolStruct{}), true)
	})
}
