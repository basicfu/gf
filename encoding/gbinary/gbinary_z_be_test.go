// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gbinary_test

import (
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_BeEncodeAndBeDecode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for k, v := range testData {
			ve := BeEncode(v)
			ve1 := BeEncodeByLength(len(ve), v)

			//t.Logf("%s:%v, encoded:%v\n", k, v, ve)
			switch v.(type) {
			case int:
				t.Assert(BeDecodeToInt(ve), v)
				t.Assert(BeDecodeToInt(ve1), v)
			case int8:
				t.Assert(BeDecodeToInt8(ve), v)
				t.Assert(BeDecodeToInt8(ve1), v)
			case int16:
				t.Assert(BeDecodeToInt16(ve), v)
				t.Assert(BeDecodeToInt16(ve1), v)
			case int32:
				t.Assert(BeDecodeToInt32(ve), v)
				t.Assert(BeDecodeToInt32(ve1), v)
			case int64:
				t.Assert(BeDecodeToInt64(ve), v)
				t.Assert(BeDecodeToInt64(ve1), v)
			case uint:
				t.Assert(BeDecodeToUint(ve), v)
				t.Assert(BeDecodeToUint(ve1), v)
			case uint8:
				t.Assert(BeDecodeToUint8(ve), v)
				t.Assert(BeDecodeToUint8(ve1), v)
			case uint16:
				t.Assert(BeDecodeToUint16(ve1), v)
				t.Assert(BeDecodeToUint16(ve), v)
			case uint32:
				t.Assert(BeDecodeToUint32(ve1), v)
				t.Assert(BeDecodeToUint32(ve), v)
			case uint64:
				t.Assert(BeDecodeToUint64(ve), v)
				t.Assert(BeDecodeToUint64(ve1), v)
			case bool:
				t.Assert(BeDecodeToBool(ve), v)
				t.Assert(BeDecodeToBool(ve1), v)
			case string:
				t.Assert(BeDecodeToString(ve), v)
				t.Assert(BeDecodeToString(ve1), v)
			case float32:
				t.Assert(BeDecodeToFloat32(ve), v)
				t.Assert(BeDecodeToFloat32(ve1), v)
			case float64:
				t.Assert(BeDecodeToFloat64(ve), v)
				t.Assert(BeDecodeToFloat64(ve1), v)
			default:
				if v == nil {
					continue
				}
				res := make([]byte, len(ve))
				err := BeDecode(ve, res)
				if err != nil {
					t.Errorf("test data: %s, %v, error:%v", k, v, err)
				}
				t.Assert(res, v)
			}
		}
	})
}

func Test_BeEncodeStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		user := User{"wenzi1", 999, "www.baidu.com"}
		ve := BeEncode(user)
		s := BeDecodeToString(ve)
		t.Assert(string(s), s)
	})
}
