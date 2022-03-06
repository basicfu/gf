// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

// go test *.go -bench=".*" -benchmem

package gcmd_test

import (
	"github.com/basicfu/gf/os/genv"
	"testing"

	"github.com/basicfu/gf/g"
	"github.com/basicfu/gf/test/gtest"
)

func Test_Default(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		Init([]string{"gf", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"}...)
		t.Assert(len(GetArgAll()), 2)
		t.Assert(GetArg(1), "path")
		t.Assert(GetArg(100, "test"), "test")
		t.Assert(GetOpt("force"), "remove")
		t.Assert(GetOpt("n"), "root")
		t.Assert(ContainsOpt("fq"), true)
		t.Assert(ContainsOpt("p"), true)
		t.Assert(ContainsOpt("none"), false)
		t.Assert(GetOpt("none", "value"), "value")
	})
	gtest.C(t, func(t *gtest.T) {
		Init([]string{"gf", "gen", "-h"}...)
		t.Assert(len(GetArgAll()), 2)
		t.Assert(GetOpt("h"), "")
		t.Assert(ContainsOpt("h"), true)
	})
}

func Test_BuildOptions(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := BuildOptions(g.MapStrStr{
			"n": "john",
		})
		t.Assert(s, "-n=john")
	})

	gtest.C(t, func(t *gtest.T) {
		s := BuildOptions(g.MapStrStr{
			"n": "john",
		}, "-test")
		t.Assert(s, "-testn=john")
	})
}

func Test_GetWithEnv(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		genv.Set("TEST", "1")
		defer genv.Remove("TEST")
		t.Assert(GetWithEnv("test"), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		genv.Set("TEST", "1")
		defer genv.Remove("TEST")
		Init("-test", "2")
		t.Assert(GetWithEnv("test"), 2)
	})
}
