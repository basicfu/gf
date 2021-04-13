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

type stringStruct1 struct {
	Name string
}

type stringStruct2 struct {
	Name string
}

func (s *stringStruct1) String() string {
	return s.Name
}

func Test_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(String(int(123)), "123")
		t.AssertEQ(String(int(-123)), "-123")
		t.AssertEQ(String(int8(123)), "123")
		t.AssertEQ(String(int8(-123)), "-123")
		t.AssertEQ(String(int16(123)), "123")
		t.AssertEQ(String(int16(-123)), "-123")
		t.AssertEQ(String(int32(123)), "123")
		t.AssertEQ(String(int32(-123)), "-123")
		t.AssertEQ(String(int64(123)), "123")
		t.AssertEQ(String(int64(-123)), "-123")
		t.AssertEQ(String(int64(1552578474888)), "1552578474888")
		t.AssertEQ(String(int64(-1552578474888)), "-1552578474888")

		t.AssertEQ(String(uint(123)), "123")
		t.AssertEQ(String(uint8(123)), "123")
		t.AssertEQ(String(uint16(123)), "123")
		t.AssertEQ(String(uint32(123)), "123")
		t.AssertEQ(String(uint64(155257847488898765)), "155257847488898765")

		t.AssertEQ(String(float32(123.456)), "123.456")
		t.AssertEQ(String(float32(-123.456)), "-123.456")
		t.AssertEQ(String(float64(1552578474888.456)), "1552578474888.456")
		t.AssertEQ(String(float64(-1552578474888.456)), "-1552578474888.456")

		t.AssertEQ(String(true), "true")
		t.AssertEQ(String(false), "false")

		t.AssertEQ(String([]byte("bytes")), "bytes")

		t.AssertEQ(String(stringStruct1{"john"}), `{"Name":"john"}`)
		t.AssertEQ(String(&stringStruct1{"john"}), "john")

		t.AssertEQ(String(stringStruct2{"john"}), `{"Name":"john"}`)
		t.AssertEQ(String(&stringStruct2{"john"}), `{"Name":"john"}`)
	})
}
