// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gmlock_test

import (
	"testing"
	"time"

	"github.com/basicfu/gf/container/garray"
	"github.com/basicfu/gf/test/gtest"
)

func Test_Locker_RLock(t *testing.T) {
	//RLock before Lock
	gtest.C(t, func(t *gtest.T) {
		key := "testRLockBeforeLock"
		array := garray.New(true)
		go func() {
			RLock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			RUnlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			Lock(key)
			array.Append(1)
			Unlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	//Lock before RLock
	gtest.C(t, func(t *gtest.T) {
		key := "testLockBeforeRLock"
		array := garray.New(true)
		go func() {
			Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			RLock(key)
			array.Append(1)
			RUnlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	//Lock before RLocks
	gtest.C(t, func(t *gtest.T) {
		key := "testLockBeforeRLocks"
		array := garray.New(true)
		go func() {
			Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			RLock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			RUnlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			RLock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			RUnlock(key)
		}()
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_Locker_TryRLock(t *testing.T) {
	//Lock before TryRLock
	gtest.C(t, func(t *gtest.T) {
		key := "testLockBeforeTryRLock"
		array := garray.New(true)
		go func() {
			Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			if TryRLock(key) {
				array.Append(1)
				RUnlock(key)
			}
		}()
		time.Sleep(150 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})

	//Lock before TryRLocks
	gtest.C(t, func(t *gtest.T) {
		key := "testLockBeforeTryRLocks"
		array := garray.New(true)
		go func() {
			Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			if TryRLock(key) {
				array.Append(1)
				RUnlock(key)
			}
		}()
		go func() {
			time.Sleep(300 * time.Millisecond)
			if TryRLock(key) {
				array.Append(1)
				RUnlock(key)
			}
		}()
		time.Sleep(150 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_Locker_RLockFunc(t *testing.T) {
	//RLockFunc before Lock
	gtest.C(t, func(t *gtest.T) {
		key := "testRLockFuncBeforeLock"
		array := garray.New(true)
		go func() {
			RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			Lock(key)
			array.Append(1)
			Unlock(key)
		}()
		time.Sleep(150 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	//Lock before RLockFunc
	gtest.C(t, func(t *gtest.T) {
		key := "testLockBeforeRLockFunc"
		array := garray.New(true)
		go func() {
			Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			RLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	//Lock before RLockFuncs
	gtest.C(t, func(t *gtest.T) {
		key := "testLockBeforeRLockFuncs"
		array := garray.New(true)
		go func() {
			Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_Locker_TryRLockFunc(t *testing.T) {
	//Lock before TryRLockFunc
	gtest.C(t, func(t *gtest.T) {
		key := "testLockBeforeTryRLockFunc"
		array := garray.New(true)
		go func() {
			Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})

	//Lock before TryRLockFuncs
	gtest.C(t, func(t *gtest.T) {
		key := "testLockBeforeTryRLockFuncs"
		array := garray.New(true)
		go func() {
			Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(300 * time.Millisecond)
			TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}
