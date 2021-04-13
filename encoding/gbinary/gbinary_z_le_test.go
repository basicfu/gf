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

func Test_LeEncodeAndLeDecode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for k, v := range testData {
			ve := LeEncode(v)
			ve1 := LeEncodeByLength(len(ve), v)

			//t.Logf("%s:%v, encoded:%v\n", k, v, ve)
			switch v.(type) {
			case int:
				t.Assert(LeDecodeToInt(ve), v)
				t.Assert(LeDecodeToInt(ve1), v)
			case int8:
				t.Assert(LeDecodeToInt8(ve), v)
				t.Assert(LeDecodeToInt8(ve1), v)
			case int16:
				t.Assert(LeDecodeToInt16(ve), v)
				t.Assert(LeDecodeToInt16(ve1), v)
			case int32:
				t.Assert(LeDecodeToInt32(ve), v)
				t.Assert(LeDecodeToInt32(ve1), v)
			case int64:
				t.Assert(LeDecodeToInt64(ve), v)
				t.Assert(LeDecodeToInt64(ve1), v)
			case uint:
				t.Assert(LeDecodeToUint(ve), v)
				t.Assert(LeDecodeToUint(ve1), v)
			case uint8:
				t.Assert(LeDecodeToUint8(ve), v)
				t.Assert(LeDecodeToUint8(ve1), v)
			case uint16:
				t.Assert(LeDecodeToUint16(ve1), v)
				t.Assert(LeDecodeToUint16(ve), v)
			case uint32:
				t.Assert(LeDecodeToUint32(ve1), v)
				t.Assert(LeDecodeToUint32(ve), v)
			case uint64:
				t.Assert(LeDecodeToUint64(ve), v)
				t.Assert(LeDecodeToUint64(ve1), v)
			case bool:
				t.Assert(LeDecodeToBool(ve), v)
				t.Assert(LeDecodeToBool(ve1), v)
			case string:
				t.Assert(LeDecodeToString(ve), v)
				t.Assert(LeDecodeToString(ve1), v)
			case float32:
				t.Assert(LeDecodeToFloat32(ve), v)
				t.Assert(LeDecodeToFloat32(ve1), v)
			case float64:
				t.Assert(LeDecodeToFloat64(ve), v)
				t.Assert(LeDecodeToFloat64(ve1), v)
			default:
				if v == nil {
					continue
				}
				res := make([]byte, len(ve))
				err := LeDecode(ve, res)
				if err != nil {
					t.Errorf("test data: %s, %v, error:%v", k, v, err)
				}
				t.Assert(res, v)
			}
		}
	})
}

func Test_LeEncodeStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		user := User{"wenzi1", 999, "www.baidu.com"}
		ve := LeEncode(user)
		s := LeDecodeToString(ve)
		t.Assert(s, s)
	})
}
