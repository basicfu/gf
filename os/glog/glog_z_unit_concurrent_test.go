// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package glog_test

import (
	"github.com/basicfu/gf/os/gfile"
	"github.com/basicfu/gf/os/gtime"
	"github.com/basicfu/gf/test/gtest"
	"github.com/basicfu/gf/text/gstr"
	"sync"
	"testing"
)

func Test_Concurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := 1000
		l := New()
		s := "@1234567890#"
		f := "test.log"
		p := gfile.TempDir(gtime.TimestampNanoStr())
		t.Assert(l.SetPath(p), nil)
		defer gfile.Remove(p)
		wg := sync.WaitGroup{}
		ch := make(chan struct{})
		for i := 0; i < c; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-ch
				l.File(f).Stdout(false).Print(s)
			}()
		}
		close(ch)
		wg.Wait()
		content := gfile.GetContents(gfile.Join(p, f))
		t.Assert(gstr.Count(content, s), c)
	})
}
