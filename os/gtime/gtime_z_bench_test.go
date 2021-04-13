// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gtime_test

import (
	"testing"
	"time"
)

func Benchmark_Timestamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Timestamp()
	}
}

func Benchmark_TimestampMilli(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TimestampMilli()
	}
}

func Benchmark_TimestampMicro(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TimestampMicro()
	}
}

func Benchmark_TimestampNano(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TimestampNano()
	}
}

func Benchmark_StrToTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StrToTime("2018-02-09T20:46:17.897Z")
	}
}

func Benchmark_StrToTime_Format(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StrToTime("2018-02-09 20:46:17.897", "Y-m-d H:i:su")
	}
}

func Benchmark_StrToTime_Layout(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StrToTimeLayout("2018-02-09T20:46:17.897Z", time.RFC3339)
	}
}

func Benchmark_ParseTimeFromContent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseTimeFromContent("2018-02-09T20:46:17.897Z")
	}
}

func Benchmark_NewFromTimeStamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewFromTimeStamp(1542674930)
	}
}

func Benchmark_Date(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Date()
	}
}

func Benchmark_Datetime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Datetime()
	}
}

func Benchmark_SetTimeZone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetTimeZone("Asia/Shanghai")
	}
}
