// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gutil_test

import (
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_ComparatorString(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorString(1, 1), 0)
		t.Assert(ComparatorString(1, 2), -1)
		t.Assert(ComparatorString(2, 1), 1)
	})
}

func Test_ComparatorInt(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorInt(1, 1), 0)
		t.Assert(ComparatorInt(1, 2), -1)
		t.Assert(ComparatorInt(2, 1), 1)
	})
}

func Test_ComparatorInt8(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorInt8(1, 1), 0)
		t.Assert(ComparatorInt8(1, 2), -1)
		t.Assert(ComparatorInt8(2, 1), 1)
	})
}

func Test_ComparatorInt16(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorInt16(1, 1), 0)
		t.Assert(ComparatorInt16(1, 2), -1)
		t.Assert(ComparatorInt16(2, 1), 1)
	})
}

func Test_ComparatorInt32(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorInt32(1, 1), 0)
		t.Assert(ComparatorInt32(1, 2), -1)
		t.Assert(ComparatorInt32(2, 1), 1)
	})
}

func Test_ComparatorInt64(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorInt64(1, 1), 0)
		t.Assert(ComparatorInt64(1, 2), -1)
		t.Assert(ComparatorInt64(2, 1), 1)
	})
}

func Test_ComparatorUint(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorUint(1, 1), 0)
		t.Assert(ComparatorUint(1, 2), -1)
		t.Assert(ComparatorUint(2, 1), 1)
	})
}

func Test_ComparatorUint8(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorUint8(1, 1), 0)
		t.Assert(ComparatorUint8(2, 6), 252)
		t.Assert(ComparatorUint8(2, 1), 1)
	})
}

func Test_ComparatorUint16(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorUint16(1, 1), 0)
		t.Assert(ComparatorUint16(1, 2), 65535)
		t.Assert(ComparatorUint16(2, 1), 1)
	})
}

func Test_ComparatorUint32(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorUint32(1, 1), 0)
		t.Assert(ComparatorUint32(-1000, 2147483640), 2147482656)
		t.Assert(ComparatorUint32(2, 1), 1)
	})
}

func Test_ComparatorUint64(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorUint64(1, 1), 0)
		t.Assert(ComparatorUint64(1, 2), -1)
		t.Assert(ComparatorUint64(2, 1), 1)
	})
}

func Test_ComparatorFloat32(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorFloat32(1, 1), 0)
		t.Assert(ComparatorFloat32(1, 2), -1)
		t.Assert(ComparatorFloat32(2, 1), 1)
	})
}

func Test_ComparatorFloat64(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorFloat64(1, 1), 0)
		t.Assert(ComparatorFloat64(1, 2), -1)
		t.Assert(ComparatorFloat64(2, 1), 1)
	})
}

func Test_ComparatorByte(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorByte(1, 1), 0)
		t.Assert(ComparatorByte(1, 2), 255)
		t.Assert(ComparatorByte(2, 1), 1)
	})
}

func Test_ComparatorRune(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(ComparatorRune(1, 1), 0)
		t.Assert(ComparatorRune(1, 2), -1)
		t.Assert(ComparatorRune(2, 1), 1)
	})
}

func Test_ComparatorTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		j := ComparatorTime("2019-06-14", "2019-06-14")
		t.Assert(j, 0)

		k := ComparatorTime("2019-06-15", "2019-06-14")
		t.Assert(k, 1)

		l := ComparatorTime("2019-06-13", "2019-06-14")
		t.Assert(l, -1)
	})
}
