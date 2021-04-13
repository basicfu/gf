// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gvalid_test

import (
	"errors"
	"github.com/basicfu/gf/frame/g"
	"github.com/basicfu/gf/util/gconv"
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_CustomRule1(t *testing.T) {
	rule := "custom"
	err := RegisterRule(rule, func(rule string, value interface{}, message string, params map[string]interface{}) error {
		pass := gconv.String(value)
		if len(pass) != 6 {
			return errors.New(message)
		}
		if params["data"] != pass {
			return errors.New(message)
		}
		return nil
	})
	gtest.Assert(err, nil)
	gtest.C(t, func(t *gtest.T) {
		err := Check("123456", rule, "custom message")
		t.Assert(err.String(), "custom message")
		err = Check("123456", rule, "custom message", g.Map{"data": "123456"})
		t.Assert(err, nil)
	})
	// Error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value string `v:"uid@custom#自定义错误"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123",
			Data:  "123456",
		}
		err := CheckStruct(st, nil)
		t.Assert(err.String(), "自定义错误")
	})
	// No error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value string `v:"uid@custom#自定义错误"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "123456",
			Data:  "123456",
		}
		err := CheckStruct(st, nil)
		t.Assert(err, nil)
	})
}

func Test_CustomRule2(t *testing.T) {
	rule := "required-map"
	err := RegisterRule(rule, func(rule string, value interface{}, message string, params map[string]interface{}) error {
		m := gconv.Map(value)
		if len(m) == 0 {
			return errors.New(message)
		}
		return nil
	})
	gtest.Assert(err, nil)
	// Check.
	gtest.C(t, func(t *gtest.T) {
		errStr := "data map should not be empty"
		t.Assert(Check(g.Map{}, rule, errStr).String(), errStr)
		t.Assert(Check(g.Map{"k": "v"}, rule, errStr).String(), nil)
	})
	// Error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value map[string]string `v:"uid@required-map#自定义错误"`
			Data  string            `p:"data"`
		}
		st := &T{
			Value: map[string]string{},
			Data:  "123456",
		}
		err := CheckStruct(st, nil)
		t.Assert(err.String(), "自定义错误")
	})
	// No error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value map[string]string `v:"uid@required-map#自定义错误"`
			Data  string            `p:"data"`
		}
		st := &T{
			Value: map[string]string{"k": "v"},
			Data:  "123456",
		}
		err := CheckStruct(st, nil)
		t.Assert(err, nil)
	})
}

func Test_CustomRule_AllowEmpty(t *testing.T) {
	rule := "allow-empty-str"
	err := RegisterRule(rule, func(rule string, value interface{}, message string, params map[string]interface{}) error {
		s := gconv.String(value)
		if len(s) == 0 || s == "gf" {
			return nil
		}
		return errors.New(message)
	})
	gtest.Assert(err, nil)
	// Check.
	gtest.C(t, func(t *gtest.T) {
		errStr := "error"
		t.Assert(Check("", rule, errStr).String(), "")
		t.Assert(Check("gf", rule, errStr).String(), "")
		t.Assert(Check("gf2", rule, errStr).String(), errStr)
	})
	// Error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value string `v:"uid@allow-empty-str#自定义错误"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "",
			Data:  "123456",
		}
		err := CheckStruct(st, nil)
		t.Assert(err.String(), "")
	})
	// No error with struct validation.
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Value string `v:"uid@allow-empty-str#自定义错误"`
			Data  string `p:"data"`
		}
		st := &T{
			Value: "john",
			Data:  "123456",
		}
		err := CheckStruct(st, nil)
		t.Assert(err.String(), "自定义错误")
	})
}
