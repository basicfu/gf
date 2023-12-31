// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

// Package garray provides most commonly used array containers which also support concurrent-safe/unsafe switch feature.
package garray

func IndexOf[T comparable](slice []T, element T) int {
	count := len(slice)
	if count == 0 {
		return -1
	}
	result := -1
	for index, v := range slice {
		if v == element {
			result = index
			break
		}
	}
	return result
}
func Contains[T comparable](slice []T, element T) bool {
	return IndexOf(slice, element) != -1
}
