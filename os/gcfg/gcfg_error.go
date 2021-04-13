// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gcfg

import (
	"github.com/basicfu/gf/os/gcmd"
)

const (
	// gERROR_PRINT_KEY is used to specify the key controlling error printing to stdout.
	// This error is designed not to be returned by functions.
	gERROR_PRINT_KEY = "gf.gcfg.errorprint"
)

// errorPrint checks whether printing error to stdout.
func errorPrint() bool {
	return gcmd.GetWithEnv(gERROR_PRINT_KEY, true).Bool()
}
