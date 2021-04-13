// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

// go test *.go -bench=".*" -count=1

package grpool_test

import (
	"testing"
)

var n = 500000

func BenchmarkGrpool2(b *testing.B) {
	b.N = n
	for i := 0; i < b.N; i++ {
		Add(increment)
	}
}

func BenchmarkGoroutine2(b *testing.B) {
	b.N = n
	for i := 0; i < b.N; i++ {
		go increment()
	}
}
