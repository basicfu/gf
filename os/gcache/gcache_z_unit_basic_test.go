// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

// go test *.go -bench=".*" -benchmem

package gcache_test

import (
	"context"
	"github.com/basicfu/gf/util/guid"
	"math"
	"testing"
	"time"

	"github.com/basicfu/gf/container/gset"
	"github.com/basicfu/gf/g"
	"github.com/basicfu/gf/os/grpool"
	"github.com/basicfu/gf/test/gtest"
)

func TestCache_GCache_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		Set(1, 11, 0)
		defer Removes(g.Slice{1, 2, 3})
		v, _ := Get(1)
		t.Assert(v, 11)
		b, _ := Contains(1)
		t.Assert(b, true)
	})
}

func TestCache_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := New()
		defer c.Close()
		t.Assert(c.Set(1, 11, 0), nil)
		v, _ := c.Get(1)
		t.Assert(v, 11)
		b, _ := c.Contains(1)
		t.Assert(b, true)
	})
}

func TestCache_GetVar(t *testing.T) {
	c := New()
	defer c.Close()
	gtest.C(t, func(t *gtest.T) {
		t.Assert(c.Set(1, 11, 0), nil)
		v, _ := c.Get(1)
		t.Assert(v, 11)
		b, _ := c.Contains(1)
		t.Assert(b, true)
	})
	gtest.C(t, func(t *gtest.T) {
		v, _ := c.GetVar(1)
		t.Assert(v.Int(), 11)
		v, _ = c.GetVar(2)
		t.Assert(v.Int(), 0)
		t.Assert(v.IsNil(), true)
		t.Assert(v.IsEmpty(), true)
	})
}

func TestCache_Set_Expire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := New()
		t.Assert(cache.Set(2, 22, 100*time.Millisecond), nil)
		v, _ := cache.Get(2)
		t.Assert(v, 22)
		time.Sleep(200 * time.Millisecond)
		v, _ = cache.Get(2)
		t.Assert(v, nil)
		time.Sleep(3 * time.Second)
		n, _ := cache.Size()
		t.Assert(n, 0)
		t.Assert(cache.Close(), nil)
	})

	gtest.C(t, func(t *gtest.T) {
		cache := New()
		t.Assert(cache.Set(1, 11, 100*time.Millisecond), nil)
		v, _ := cache.Get(1)
		t.Assert(v, 11)
		time.Sleep(200 * time.Millisecond)
		v, _ = cache.Get(1)
		t.Assert(v, nil)
	})
}

func TestCache_Update_GetExpire(t *testing.T) {
	// gcache
	gtest.C(t, func(t *gtest.T) {
		key := guid.S()
		Set(key, 11, 3*time.Second)
		expire1, _ := GetExpire(key)
		Update(key, 12)
		expire2, _ := GetExpire(key)
		v, _ := GetVar(key)
		t.Assert(v, 12)
		t.Assert(math.Ceil(expire1.Seconds()), math.Ceil(expire2.Seconds()))
	})
	// gcache.Cache
	gtest.C(t, func(t *gtest.T) {
		cache := New()
		cache.Set(1, 11, 3*time.Second)
		expire1, _ := cache.GetExpire(1)
		cache.Update(1, 12)
		expire2, _ := cache.GetExpire(1)
		v, _ := cache.GetVar(1)
		t.Assert(v, 12)
		t.Assert(math.Ceil(expire1.Seconds()), math.Ceil(expire2.Seconds()))
	})
}

func TestCache_UpdateExpire(t *testing.T) {
	// gcache
	gtest.C(t, func(t *gtest.T) {
		key := guid.S()
		Set(key, 11, 3*time.Second)
		defer Remove(key)
		oldExpire, _ := GetExpire(key)
		newExpire := 10 * time.Second
		UpdateExpire(key, newExpire)
		e, _ := GetExpire(key)
		t.AssertNE(e, oldExpire)
		e, _ = GetExpire(key)
		t.Assert(math.Ceil(e.Seconds()), 10)
	})
	// gcache.Cache
	gtest.C(t, func(t *gtest.T) {
		cache := New()
		cache.Set(1, 11, 3*time.Second)
		oldExpire, _ := cache.GetExpire(1)
		newExpire := 10 * time.Second
		cache.UpdateExpire(1, newExpire)
		e, _ := cache.GetExpire(1)
		t.AssertNE(e, oldExpire)

		e, _ = cache.GetExpire(1)
		t.Assert(math.Ceil(e.Seconds()), 10)
	})
}

func TestCache_Keys_Values(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := New()
		for i := 0; i < 10; i++ {
			t.Assert(c.Set(i, i*10, 0), nil)
		}
		var (
			keys, _   = c.Keys()
			values, _ = c.Values()
		)
		t.Assert(len(keys), 10)
		t.Assert(len(values), 10)
		t.AssertIN(0, keys)
		t.AssertIN(90, values)
	})
}

func TestCache_LRU(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := New(2)
		for i := 0; i < 10; i++ {
			cache.Set(i, i, 0)
		}
		n, _ := cache.Size()
		t.Assert(n, 10)
		v, _ := cache.Get(6)
		t.Assert(v, 6)
		time.Sleep(4 * time.Second)
		n, _ = cache.Size()
		t.Assert(n, 2)
		v, _ = cache.Get(6)
		t.Assert(v, 6)
		v, _ = cache.Get(1)
		t.Assert(v, nil)
		t.Assert(cache.Close(), nil)
	})
}

func TestCache_LRU_expire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := New(2)
		t.Assert(cache.Set(1, nil, 1000), nil)
		n, _ := cache.Size()
		t.Assert(n, 1)
		v, _ := cache.Get(1)

		t.Assert(v, nil)
	})
}

func TestCache_SetIfNotExist(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := New()
		cache.SetIfNotExist(1, 11, 0)
		v, _ := cache.Get(1)
		t.Assert(v, 11)
		cache.SetIfNotExist(1, 22, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)
		cache.SetIfNotExist(2, 22, 0)
		v, _ = cache.Get(2)
		t.Assert(v, 22)

		Removes(g.Slice{1, 2, 3})
		SetIfNotExist(1, 11, 0)
		v, _ = Get(1)
		t.Assert(v, 11)
		SetIfNotExist(1, 22, 0)
		v, _ = Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_Sets(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := New()
		cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		v, _ := cache.Get(1)
		t.Assert(v, 11)

		Removes(g.Slice{1, 2, 3})
		Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_GetOrSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := New()
		cache.GetOrSet(1, 11, 0)
		v, _ := cache.Get(1)
		t.Assert(v, 11)
		cache.GetOrSet(1, 111, 0)

		v, _ = cache.Get(1)
		t.Assert(v, 11)
		Removes(g.Slice{1, 2, 3})
		GetOrSet(1, 11, 0)

		v, _ = cache.Get(1)
		t.Assert(v, 11)

		GetOrSet(1, 111, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_GetOrSetFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := New()
		cache.GetOrSetFunc(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		v, _ := cache.Get(1)
		t.Assert(v, 11)

		cache.GetOrSetFunc(1, func() (interface{}, error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)

		Removes(g.Slice{1, 2, 3})

		GetOrSetFunc(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)

		GetOrSetFunc(1, func() (interface{}, error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_GetOrSetFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := New()
		cache.GetOrSetFuncLock(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		v, _ := cache.Get(1)
		t.Assert(v, 11)

		cache.GetOrSetFuncLock(1, func() (interface{}, error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)

		Removes(g.Slice{1, 2, 3})
		GetOrSetFuncLock(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)

		GetOrSetFuncLock(1, func() (interface{}, error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := New()
		cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		cache.Clear()
		n, _ := cache.Size()
		t.Assert(n, 0)
	})
}

func TestCache_SetConcurrency(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := New()
		pool := grpool.New(4)
		go func() {
			for {
				pool.Add(func() {
					cache.SetIfNotExist(1, 11, 10)
				})
			}
		}()
		select {
		case <-time.After(2 * time.Second):
			//t.Log("first part end")
		}

		go func() {
			for {
				pool.Add(func() {
					cache.SetIfNotExist(1, nil, 10)
				})
			}
		}()
		select {
		case <-time.After(2 * time.Second):
			//t.Log("second part end")
		}
	})
}

func TestCache_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		{
			cache := New()
			cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
			b, _ := cache.Contains(1)
			t.Assert(b, true)
			v, _ := cache.Get(1)
			t.Assert(v, 11)
			data, _ := cache.Data()
			t.Assert(data[1], 11)
			t.Assert(data[2], 22)
			t.Assert(data[3], nil)
			n, _ := cache.Size()
			t.Assert(n, 2)
			keys, _ := cache.Keys()
			t.Assert(gset.NewFrom(g.Slice{1, 2}).Equal(gset.NewFrom(keys)), true)
			keyStrs, _ := cache.KeyStrings()
			t.Assert(gset.NewFrom(g.Slice{"1", "2"}).Equal(gset.NewFrom(keyStrs)), true)
			values, _ := cache.Values()
			t.Assert(gset.NewFrom(g.Slice{11, 22}).Equal(gset.NewFrom(values)), true)
			removeData1, _ := cache.Remove(1)
			t.Assert(removeData1, 11)
			n, _ = cache.Size()
			t.Assert(n, 1)
			cache.Removes(g.Slice{2})
			n, _ = cache.Size()
			t.Assert(n, 0)
		}

		Remove(g.Slice{1, 2, 3}...)
		{
			Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
			b, _ := Contains(1)
			t.Assert(b, true)
			v, _ := Get(1)
			t.Assert(v, 11)
			data, _ := Data()
			t.Assert(data[1], 11)
			t.Assert(data[2], 22)
			t.Assert(data[3], nil)
			n, _ := Size()
			t.Assert(n, 2)
			keys, _ := Keys()
			t.Assert(gset.NewFrom(g.Slice{1, 2}).Equal(gset.NewFrom(keys)), true)
			keyStrs, _ := KeyStrings()
			t.Assert(gset.NewFrom(g.Slice{"1", "2"}).Equal(gset.NewFrom(keyStrs)), true)
			values, _ := Values()
			t.Assert(gset.NewFrom(g.Slice{11, 22}).Equal(gset.NewFrom(values)), true)
			removeData1, _ := Remove(1)
			t.Assert(removeData1, 11)
			n, _ = Size()
			t.Assert(n, 1)
			Removes(g.Slice{2})
			n, _ = Size()
			t.Assert(n, 0)
		}
	})
}

func TestCache_Ctx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := New()
		cache.Ctx(context.Background()).Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		b, _ := cache.Contains(1)
		t.Assert(b, true)
		v, _ := cache.Get(1)
		t.Assert(v, 11)
	})
}
