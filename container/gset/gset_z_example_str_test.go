// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/basicfu/gf.

package gset_test

import (
	"fmt"
	"github.com/basicfu/gf/g"
)

func ExampleStrSet_Contains() {
	var set StrSet
	set.Add("a")
	fmt.Println(set.Contains("a"))
	fmt.Println(set.Contains("A"))
	fmt.Println(set.ContainsI("A"))

	// Output:
	// true
	// false
	// true
}

func ExampleStrSet_Walk() {
	var (
		set    StrSet
		names  = g.SliceStr{"user", "user_detail"}
		prefix = "gf_"
	)
	set.Add(names...)
	// Add prefix for given table names.
	set.Walk(func(item string) string {
		return prefix + item
	})
	fmt.Println(set.Slice())

	// May Output:
	// [gf_user gf_user_detail]
}
