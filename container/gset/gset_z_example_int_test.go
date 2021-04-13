// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/basicfu/gf.

package gset_test

import (
	"fmt"
)

func ExampleIntSet_Contains() {
	var set IntSet
	set.Add(1)
	fmt.Println(set.Contains(1))
	fmt.Println(set.Contains(2))

	// Output:
	// true
	// false
}
