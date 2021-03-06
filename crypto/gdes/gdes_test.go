// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gdes_test

import (
	"encoding/hex"
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

var (
	errKey     = []byte("1111111111111234123456789")
	errIv      = []byte("123456789")
	errPadding = 5
)

func TestDesECB(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		key := []byte("11111111")
		text := []byte("12345678")
		padding := NOPADDING
		result := "858b176da8b12503"
		// encrypt test
		cipherText, err := EncryptECB(text, key, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(hex.EncodeToString(cipherText), result)
		// decrypt test
		clearText, err := DecryptECB(cipherText, key, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(string(clearText), "12345678")

		// encrypt err test. when throw exception,the err is not equal nil and the string is nil
		errEncrypt, err := EncryptECB(text, key, errPadding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
		errEncrypt, err = EncryptECB(text, errKey, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
		// err decrypt test.
		errDecrypt, err := DecryptECB(cipherText, errKey, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errDecrypt, nil)
		errDecrypt, err = DecryptECB(cipherText, key, errPadding)
		t.AssertNE(err, nil)
		t.AssertEQ(errDecrypt, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		key := []byte("11111111")
		text := []byte("12345678")
		padding := PKCS5PADDING
		errPadding := 5
		result := "858b176da8b12503ad6a88b4fa37833d"
		cipherText, err := EncryptECB(text, key, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(hex.EncodeToString(cipherText), result)
		// decrypt test
		clearText, err := DecryptECB(cipherText, key, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(string(clearText), "12345678")

		// err test
		errEncrypt, err := EncryptECB(text, key, errPadding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
		errDecrypt, err := DecryptECB(cipherText, errKey, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errDecrypt, nil)
	})
}

func Test3DesECB(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		key := []byte("1111111111111234")
		text := []byte("1234567812345678")
		padding := NOPADDING
		result := "a23ee24b98c26263a23ee24b98c26263"
		// encrypt test
		cipherText, err := EncryptECBTriple(text, key, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(hex.EncodeToString(cipherText), result)
		// decrypt test
		clearText, err := DecryptECBTriple(cipherText, key, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(string(clearText), "1234567812345678")
		// err test
		errEncrypt, err := EncryptECB(text, key, errPadding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		key := []byte("111111111111123412345678")
		text := []byte("123456789")
		padding := PKCS5PADDING
		errPadding := 5
		result := "37989b1effc07a6d00ff89a7d052e79f"
		// encrypt test
		cipherText, err := EncryptECBTriple(text, key, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(hex.EncodeToString(cipherText), result)
		// decrypt test
		clearText, err := DecryptECBTriple(cipherText, key, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(string(clearText), "123456789")
		// err test, when key is err, but text and padding is right
		errEncrypt, err := EncryptECBTriple(text, errKey, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
		// when padding is err,but key and text is right
		errEncrypt, err = EncryptECBTriple(text, key, errPadding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
		// decrypt err test,when key is err
		errEncrypt, err = DecryptECBTriple(text, errKey, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
	})
}

func TestDesCBC(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		key := []byte("11111111")
		text := []byte("1234567812345678")
		padding := NOPADDING
		iv := []byte("12345678")
		result := "40826a5800608c87585ca7c9efabee47"
		// encrypt test
		cipherText, err := EncryptCBC(text, key, iv, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(hex.EncodeToString(cipherText), result)
		// decrypt test
		clearText, err := DecryptCBC(cipherText, key, iv, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(string(clearText), "1234567812345678")
		// encrypt err test.
		errEncrypt, err := EncryptCBC(text, errKey, iv, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
		// the iv is err
		errEncrypt, err = EncryptCBC(text, key, errIv, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
		// the padding is err
		errEncrypt, err = EncryptCBC(text, key, iv, errPadding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
		// decrypt err test. the key is err
		errDecrypt, err := DecryptCBC(cipherText, errKey, iv, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errDecrypt, nil)
		// the iv is err
		errDecrypt, err = DecryptCBC(cipherText, key, errIv, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errDecrypt, nil)
		// the padding is err
		errDecrypt, err = DecryptCBC(cipherText, key, iv, errPadding)
		t.AssertNE(err, nil)
		t.AssertEQ(errDecrypt, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		key := []byte("11111111")
		text := []byte("12345678")
		padding := PKCS5PADDING
		iv := []byte("12345678")
		result := "40826a5800608c87100a25d86ac7c52c"
		// encrypt test
		cipherText, err := EncryptCBC(text, key, iv, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(hex.EncodeToString(cipherText), result)
		// decrypt test
		clearText, err := DecryptCBC(cipherText, key, iv, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(string(clearText), "12345678")
		// err test
		errEncrypt, err := EncryptCBC(text, key, errIv, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
	})
}

func Test3DesCBC(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		key := []byte("1111111112345678")
		text := []byte("1234567812345678")
		padding := NOPADDING
		iv := []byte("12345678")
		result := "bfde1394e265d5f738d5cab170c77c88"
		// encrypt test
		cipherText, err := EncryptCBCTriple(text, key, iv, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(hex.EncodeToString(cipherText), result)
		// decrypt test
		clearText, err := DecryptCBCTriple(cipherText, key, iv, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(string(clearText), "1234567812345678")
		// encrypt err test
		errEncrypt, err := EncryptCBCTriple(text, errKey, iv, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
		// the iv is err
		errEncrypt, err = EncryptCBCTriple(text, key, errIv, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
		// the padding is err
		errEncrypt, err = EncryptCBCTriple(text, key, iv, errPadding)
		t.AssertNE(err, nil)
		t.AssertEQ(errEncrypt, nil)
		// decrypt err test
		errDecrypt, err := DecryptCBCTriple(cipherText, errKey, iv, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errDecrypt, nil)
		// the iv is err
		errDecrypt, err = DecryptCBCTriple(cipherText, key, errIv, padding)
		t.AssertNE(err, nil)
		t.AssertEQ(errDecrypt, nil)
		// the padding is err
		errDecrypt, err = DecryptCBCTriple(cipherText, key, iv, errPadding)
		t.AssertNE(err, nil)
		t.AssertEQ(errDecrypt, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		key := []byte("111111111234567812345678")
		text := []byte("12345678")
		padding := PKCS5PADDING
		iv := []byte("12345678")
		result := "40826a5800608c87100a25d86ac7c52c"
		// encrypt test
		cipherText, err := EncryptCBCTriple(text, key, iv, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(hex.EncodeToString(cipherText), result)
		// decrypt test
		clearText, err := DecryptCBCTriple(cipherText, key, iv, padding)
		t.AssertEQ(err, nil)
		t.AssertEQ(string(clearText), "12345678")
	})

}
