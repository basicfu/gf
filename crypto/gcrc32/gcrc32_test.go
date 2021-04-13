// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

// go test *.go -bench=".*"

package gcrc32_test

import (
	"testing"

	"github.com/basicfu/gf/crypto/gmd5"
	"github.com/basicfu/gf/test/gtest"
)

func TestEncrypt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "pibigstar"
		result := 693191136
		encrypt1 := Encrypt(s)
		encrypt2 := Encrypt([]byte(s))
		t.AssertEQ(int(encrypt1), result)
		t.AssertEQ(int(encrypt2), result)

		strmd5, _ := gmd5.Encrypt(s)
		test1 := Encrypt(strmd5)
		test2 := Encrypt([]byte(strmd5))
		t.AssertEQ(test2, test1)
	})
}
