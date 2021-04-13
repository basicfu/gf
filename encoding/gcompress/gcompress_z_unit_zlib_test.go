// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gcompress_test

import (
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_Zlib_UnZlib(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		src := "hello, world\n"
		dst := []byte{120, 156, 202, 72, 205, 201, 201, 215, 81, 40, 207, 47, 202, 73, 225, 2, 4, 0, 0, 255, 255, 33, 231, 4, 147}
		data, _ := Zlib([]byte(src))
		t.Assert(data, dst)

		data, _ = UnZlib(dst)
		t.Assert(data, []byte(src))

		data, _ = Zlib(nil)
		t.Assert(data, nil)
		data, _ = UnZlib(nil)
		t.Assert(data, nil)

		data, _ = UnZlib(dst[1:])
		t.Assert(data, nil)
	})
}
