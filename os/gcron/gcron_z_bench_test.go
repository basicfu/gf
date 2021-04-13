// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gcron_test

import (
	"testing"
)

func Benchmark_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add("1 1 1 1 1 1", func() {

		})
	}
}
