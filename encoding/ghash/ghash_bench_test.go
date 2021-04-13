// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

// go test *.go -bench=".*"

package ghash_test

import (
	"testing"
)

var (
	str = []byte("This is the test string for hash.")
)

func BenchmarkBKDRHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BKDRHash(str)
	}
}

func BenchmarkBKDRHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BKDRHash64(str)
	}
}

func BenchmarkSDBMHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SDBMHash(str)
	}
}

func BenchmarkSDBMHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SDBMHash64(str)
	}
}

func BenchmarkRSHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RSHash(str)
	}
}

func BenchmarkSRSHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RSHash64(str)
	}
}

func BenchmarkJSHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JSHash(str)
	}
}

func BenchmarkJSHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JSHash64(str)
	}
}

func BenchmarkPJWHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PJWHash(str)
	}
}

func BenchmarkPJWHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PJWHash64(str)
	}
}

func BenchmarkELFHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ELFHash(str)
	}
}

func BenchmarkELFHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ELFHash64(str)
	}
}

func BenchmarkDJBHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DJBHash(str)
	}
}

func BenchmarkDJBHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DJBHash64(str)
	}
}

func BenchmarkAPHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		APHash(str)
	}
}

func BenchmarkAPHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		APHash64(str)
	}
}
