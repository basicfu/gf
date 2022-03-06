// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package g

import (
	"github.com/basicfu/gf/container/gvar"
)

// NewVar returns a gvar.Var.
func NewVar(i interface{}, safe ...bool) *Var {
	return gvar.New(i, safe...)
}
