// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gconv

import "github.com/shopspring/decimal"

// 目前接受string，可使用泛型判断
func Decimal(str string) decimal.Decimal {
	amount, err := decimal.NewFromString(str)
	if err != nil {
		return decimal.NewFromInt(0)
	}
	return amount
}
