// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package g

import (
	"fmt"
	"runtime/debug"

	"github.com/basicfu/gf/decimal"
	"github.com/basicfu/gf/errors/gerror"
)

// 全局错误不会捕捉，也可以加入全局捕捉，嵌套一层catch
var Go = func(handler func(), catch ...func(err error)) {
	go TryBlock(handler, catch...)
}

var Decimal = decimal.New

// 仅接受分为单位，入参应为整数分，输出参数固定2为小数金额，虽为四舍五入但实际不存在四舍五入的情况
var DecimalCentStr = func(amount any) string {
	return decimal.New(amount).Div(decimal.New(100)).StringFixed(2)
}

// Frequently-used map type alias.
type Map = map[string]interface{}
type MapAnyAny = map[interface{}]interface{}
type MapAnyStr = map[interface{}]string
type MapAnyInt = map[interface{}]int
type MapStrAny = map[string]interface{}
type MapStrStr = map[string]string
type MapStrInt = map[string]int
type MapIntAny = map[int]interface{}
type MapIntStr = map[int]string
type MapIntInt = map[int]int
type MapAnyBool = map[interface{}]bool
type MapStrBool = map[string]bool
type MapIntBool = map[int]bool

// Frequently-used slice type alias.
type List = []Map
type ListAnyAny = []MapAnyAny
type ListAnyStr = []MapAnyStr
type ListAnyInt = []MapAnyInt
type ListStrAny = []MapStrAny
type ListStrStr = []MapStrStr
type ListStrInt = []MapStrInt
type ListIntAny = []MapIntAny
type ListIntStr = []MapIntStr
type ListIntInt = []MapIntInt
type ListAnyBool = []MapAnyBool
type ListStrBool = []MapStrBool
type ListIntBool = []MapIntBool

// Frequently-used slice type alias.
type Slice = []interface{}
type SliceAny = []interface{}
type SliceStr = []string
type SliceInt = []int

// Array is alias of Slice.
type Array = []interface{}
type ArrayAny = []interface{}
type ArrayStr = []string
type ArrayInt = []int

func Try(catch ...func(err error)) {
	if err := recover(); err != nil {
		if len(catch) == 0 {
			return
		}
		if v, ok := err.(error); ok {
			catch[0](v)
		} else {
			catch[0](gerror.Newf(`%+v`, err))
		}
	}
}
func TryBlock(try func(), catch ...func(err error)) {
	defer func() {
		if err := recover(); err != nil {
			if len(catch) == 0 {
				fmt.Println("go error:", gerror.Newf(`%+v`, err), string(debug.Stack()))
				return
			}
			if v, ok := err.(error); ok {
				catch[0](v)
			} else {
				catch[0](gerror.Newf(`%+v`, err))
			}
		}
	}()
	try()
}

// TryCatch implements try...catch... logistics using internal panic...recover.
// It automatically calls function <catch> if any exception occurs ans passes the exception as an error.
func TryCatch(try func(), catch ...func(exception error)) {
	defer func() {
		if exception := recover(); exception != nil && len(catch) > 0 {
			if err, ok := exception.(error); ok {
				catch[0](err)
			} else {
				catch[0](fmt.Errorf(`%v`, exception))
			}
		}
	}()
	try()
}
