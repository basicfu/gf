// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

// go test *.go -bench=".*"

package grand_test

import (
	"github.com/basicfu/gf/text/gstr"
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_Intn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 1000000; i++ {
			n := Intn(100)
			t.AssertLT(n, 100)
			t.AssertGE(n, 0)
		}
		for i := 0; i < 1000000; i++ {
			n := Intn(-100)
			t.AssertLE(n, 0)
			t.Assert(n, -100)
		}
	})
}

func Test_Meet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(Meet(100, 100), true)
		}
		for i := 0; i < 100; i++ {
			t.Assert(Meet(0, 100), false)
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(Meet(50, 100), []bool{true, false})
		}
	})
}

func Test_MeetProb(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(MeetProb(1), true)
		}
		for i := 0; i < 100; i++ {
			t.Assert(MeetProb(0), false)
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(MeetProb(0.5), []bool{true, false})
		}
	})
}

func Test_N(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(N(1, 1), 1)
		}
		for i := 0; i < 100; i++ {
			t.Assert(N(0, 0), 0)
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(N(1, 2), []int{1, 2})
		}
	})
}

func Test_Rand(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(N(1, 1), 1)
		}
		for i := 0; i < 100; i++ {
			t.Assert(N(0, 0), 0)
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(N(1, 2), []int{1, 2})
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(N(-1, 2), []int{-1, 0, 1, 2})
		}
	})
}

func Test_S(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(S(5)), 5)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(S(5, true)), 5)
		}
	})
}

func Test_B(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			b := B(5)
			t.Assert(len(b), 5)
			t.AssertNE(b, make([]byte, 5))
		}
	})
}

func Test_Str(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(S(5)), 5)
		}
	})
}

func Test_RandStr(t *testing.T) {
	str := "我爱GoFrame"
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 10; i++ {
			s := Str(str, 100000)
			t.Assert(gstr.Contains(s, "我"), true)
			t.Assert(gstr.Contains(s, "爱"), true)
			t.Assert(gstr.Contains(s, "G"), true)
			t.Assert(gstr.Contains(s, "o"), true)
			t.Assert(gstr.Contains(s, "F"), true)
			t.Assert(gstr.Contains(s, "r"), true)
			t.Assert(gstr.Contains(s, "a"), true)
			t.Assert(gstr.Contains(s, "m"), true)
			t.Assert(gstr.Contains(s, "e"), true)
			t.Assert(gstr.Contains(s, "w"), false)
		}
	})
}

func Test_Digits(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(Digits(5)), 5)
		}
	})
}

func Test_RandDigits(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(Digits(5)), 5)
		}
	})
}

func Test_Letters(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(Letters(5)), 5)
		}
	})
}

func Test_RandLetters(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(Letters(5)), 5)
		}
	})
}

func Test_Perm(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			t.AssertIN(Perm(5), []int{0, 1, 2, 3, 4})
		}
	})
}
