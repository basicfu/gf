// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

// go test *.go -bench=".*"

package gstr_test

import (
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_Pos(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(Pos(s1, "ab"), 0)
		t.Assert(Pos(s1, "ab", 2), 7)
		t.Assert(Pos(s1, "abd", 0), -1)
		t.Assert(Pos(s1, "e", -4), 11)
	})
	gtest.C(t, func(t *gtest.T) {
		s1 := "我爱China very much"
		t.Assert(Pos(s1, "爱"), 3)
		t.Assert(Pos(s1, "C"), 6)
		t.Assert(Pos(s1, "China"), 6)
	})
}

func Test_PosRune(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(PosRune(s1, "ab"), 0)
		t.Assert(PosRune(s1, "ab", 2), 7)
		t.Assert(PosRune(s1, "abd", 0), -1)
		t.Assert(PosRune(s1, "e", -4), 11)
	})
	gtest.C(t, func(t *gtest.T) {
		s1 := "我爱China very much"
		t.Assert(PosRune(s1, "爱"), 1)
		t.Assert(PosRune(s1, "C"), 2)
		t.Assert(PosRune(s1, "China"), 2)
	})
}

func Test_PosI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(PosI(s1, "zz"), -1)
		t.Assert(PosI(s1, "ab"), 0)
		t.Assert(PosI(s1, "ef", 2), 4)
		t.Assert(PosI(s1, "abd", 0), -1)
		t.Assert(PosI(s1, "E", -4), 11)
	})
	gtest.C(t, func(t *gtest.T) {
		s1 := "我爱China very much"
		t.Assert(PosI(s1, "爱"), 3)
		t.Assert(PosI(s1, "c"), 6)
		t.Assert(PosI(s1, "china"), 6)
	})
}

func Test_PosIRune(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(PosIRune(s1, "zz"), -1)
		t.Assert(PosIRune(s1, "ab"), 0)
		t.Assert(PosIRune(s1, "ef", 2), 4)
		t.Assert(PosIRune(s1, "abd", 0), -1)
		t.Assert(PosIRune(s1, "E", -4), 11)
	})
	gtest.C(t, func(t *gtest.T) {
		s1 := "我爱China very much"
		t.Assert(PosIRune(s1, "爱"), 1)
		t.Assert(PosIRune(s1, "c"), 2)
		t.Assert(PosIRune(s1, "china"), 2)
	})
}

func Test_PosR(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(PosR(s1, "zz"), -1)
		t.Assert(PosR(s1, "ab"), 7)
		t.Assert(PosR(s2, "ab", -2), 0)
		t.Assert(PosR(s1, "ef"), 11)
		t.Assert(PosR(s1, "abd", 0), -1)
		t.Assert(PosR(s1, "e", -4), -1)
	})
	gtest.C(t, func(t *gtest.T) {
		s1 := "我爱China very much"
		t.Assert(PosR(s1, "爱"), 3)
		t.Assert(PosR(s1, "C"), 6)
		t.Assert(PosR(s1, "China"), 6)
	})
}

func Test_PosRRune(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(PosRRune(s1, "zz"), -1)
		t.Assert(PosRRune(s1, "ab"), 7)
		t.Assert(PosRRune(s2, "ab", -2), 0)
		t.Assert(PosRRune(s1, "ef"), 11)
		t.Assert(PosRRune(s1, "abd", 0), -1)
		t.Assert(PosRRune(s1, "e", -4), -1)
	})
	gtest.C(t, func(t *gtest.T) {
		s1 := "我爱China very much"
		t.Assert(PosRRune(s1, "爱"), 1)
		t.Assert(PosRRune(s1, "C"), 2)
		t.Assert(PosRRune(s1, "China"), 2)
	})
}

func Test_PosRI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(PosRI(s1, "zz"), -1)
		t.Assert(PosRI(s1, "AB"), 7)
		t.Assert(PosRI(s2, "AB", -2), 0)
		t.Assert(PosRI(s1, "EF"), 11)
		t.Assert(PosRI(s1, "abd", 0), -1)
		t.Assert(PosRI(s1, "e", -5), 4)
	})
	gtest.C(t, func(t *gtest.T) {
		s1 := "我爱China very much"
		t.Assert(PosRI(s1, "爱"), 3)
		t.Assert(PosRI(s1, "C"), 19)
		t.Assert(PosRI(s1, "China"), 6)
	})
}

func Test_PosRIRune(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(PosRIRune(s1, "zz"), -1)
		t.Assert(PosRIRune(s1, "AB"), 7)
		t.Assert(PosRIRune(s2, "AB", -2), 0)
		t.Assert(PosRIRune(s1, "EF"), 11)
		t.Assert(PosRIRune(s1, "abd", 0), -1)
		t.Assert(PosRIRune(s1, "e", -5), 4)
	})
	gtest.C(t, func(t *gtest.T) {
		s1 := "我爱China very much"
		t.Assert(PosRIRune(s1, "爱"), 1)
		t.Assert(PosRIRune(s1, "C"), 15)
		t.Assert(PosRIRune(s1, "China"), 2)
	})
}
