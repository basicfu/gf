// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

// go test *.go -bench=".*"

package gmode_test

import (
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_AutoCheckSourceCodes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(IsDevelop(), true)
	})
}

func Test_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		oldMode := Mode()
		defer Set(oldMode)
		SetDevelop()
		t.Assert(IsDevelop(), true)
		t.Assert(IsTesting(), false)
		t.Assert(IsStaging(), false)
		t.Assert(IsProduct(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		oldMode := Mode()
		defer Set(oldMode)
		SetTesting()
		t.Assert(IsDevelop(), false)
		t.Assert(IsTesting(), true)
		t.Assert(IsStaging(), false)
		t.Assert(IsProduct(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		oldMode := Mode()
		defer Set(oldMode)
		SetStaging()
		t.Assert(IsDevelop(), false)
		t.Assert(IsTesting(), false)
		t.Assert(IsStaging(), true)
		t.Assert(IsProduct(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		oldMode := Mode()
		defer Set(oldMode)
		SetProduct()
		t.Assert(IsDevelop(), false)
		t.Assert(IsTesting(), false)
		t.Assert(IsStaging(), false)
		t.Assert(IsProduct(), true)
	})
}
