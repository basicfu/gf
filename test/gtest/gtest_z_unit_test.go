// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gtest_test

import (
	"testing"
)

func TestC(t *testing.T) {
	C(t, func(t *T) {
		t.Assert(1, 1)
		t.AssertNE(1, 0)
		t.AssertEQ(float32(123.456), float32(123.456))
		t.AssertEQ(float64(123.456), float64(123.456))
	})
}

func TestCase(t *testing.T) {
	C(t, func(t *T) {
		t.Assert(1, 1)
		t.AssertNE(1, 0)
		t.AssertEQ(float32(123.456), float32(123.456))
		t.AssertEQ(float64(123.456), float64(123.456))
	})
}
