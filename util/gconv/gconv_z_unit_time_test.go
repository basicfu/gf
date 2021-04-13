// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gconv_test

import (
	"github.com/basicfu/gf/frame/g"
	"testing"
	"time"

	"github.com/basicfu/gf/os/gtime"
	"github.com/basicfu/gf/test/gtest"
)

func Test_Time(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "2011-10-10 01:02:03.456"
		t.AssertEQ(GTime(s), gtime.NewFromStr(s))
		t.AssertEQ(Time(s), gtime.NewFromStr(s).Time)
		t.AssertEQ(Duration(100), 100*time.Nanosecond)
	})
	gtest.C(t, func(t *gtest.T) {
		s := "01:02:03.456"
		t.AssertEQ(GTime(s).Hour(), 1)
		t.AssertEQ(GTime(s).Minute(), 2)
		t.AssertEQ(GTime(s).Second(), 3)
		t.AssertEQ(GTime(s), gtime.NewFromStr(s))
		t.AssertEQ(Time(s), gtime.NewFromStr(s).Time)
	})
	gtest.C(t, func(t *gtest.T) {
		s := "0000-01-01 01:02:03"
		t.AssertEQ(GTime(s).Year(), 0)
		t.AssertEQ(GTime(s).Month(), 1)
		t.AssertEQ(GTime(s).Day(), 1)
		t.AssertEQ(GTime(s).Hour(), 1)
		t.AssertEQ(GTime(s).Minute(), 2)
		t.AssertEQ(GTime(s).Second(), 3)
		t.AssertEQ(GTime(s), gtime.NewFromStr(s))
		t.AssertEQ(Time(s), gtime.NewFromStr(s).Time)
	})
}

func Test_Time_Slice_Attribute(t *testing.T) {
	type SelectReq struct {
		Arr []*gtime.Time
		One *gtime.Time
	}
	gtest.C(t, func(t *gtest.T) {
		var s *SelectReq
		err := Struct(g.Map{
			"arr": g.Slice{"2021-01-12 12:34:56", "2021-01-12 12:34:57"},
			"one": "2021-01-12 12:34:58",
		}, &s)
		t.Assert(err, nil)
		t.Assert(s.One, "2021-01-12 12:34:58")
		t.Assert(s.Arr[0], "2021-01-12 12:34:56")
		t.Assert(s.Arr[1], "2021-01-12 12:34:57")
	})
}
