// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gfile

import (
	"github.com/basicfu/gf/os/gcmd"
	"time"
)

const (
	// Default expire time for file content caching in seconds.
	gDEFAULT_CACHE_EXPIRE = time.Minute
)

var (
	// Default expire time for file content caching.
	cacheExpire = gcmd.GetWithEnv("gf.gfile.cache", gDEFAULT_CACHE_EXPIRE).Duration()

	// internalCache is the memory cache for internal usage.
)

// cacheKey produces the cache key for gcache.
func cacheKey(path string) string {
	return "gf.gfile.cache:" + path
}
