// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gins_test

import (
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_SetGet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		Set("test-user", 1)
		t.Assert(Get("test-user"), 1)
		t.Assert(Get("none-exists"), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(GetOrSet("test-1", 1), 1)
		t.Assert(Get("test-1"), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(GetOrSetFunc("test-2", func() interface{} {
			return 2
		}), 2)
		t.Assert(Get("test-2"), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(GetOrSetFuncLock("test-3", func() interface{} {
			return 3
		}), 3)
		t.Assert(Get("test-3"), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(SetIfNotExist("test-4", 4), true)
		t.Assert(Get("test-4"), 4)
		t.Assert(SetIfNotExist("test-4", 5), false)
		t.Assert(Get("test-4"), 4)
	})
}
