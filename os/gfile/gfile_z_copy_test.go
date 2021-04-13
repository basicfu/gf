// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gfile_test

import (
	"github.com/basicfu/gf/os/gtime"
	"github.com/basicfu/gf/test/gtest"
	"testing"
)

func Test_Copy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			paths  = "/testfile_copyfile1.txt"
			topath = "/testfile_copyfile2.txt"
		)

		createTestFile(paths, "")
		defer delTestFiles(paths)

		t.Assert(Copy(testpath()+paths, testpath()+topath), nil)
		defer delTestFiles(topath)

		t.Assert(IsFile(testpath()+topath), true)
		t.AssertNE(Copy("", ""), nil)
	})
}

func Test_CopyFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			paths  = "/testfile_copyfile1.txt"
			topath = "/testfile_copyfile2.txt"
		)

		createTestFile(paths, "")
		defer delTestFiles(paths)

		t.Assert(CopyFile(testpath()+paths, testpath()+topath), nil)
		defer delTestFiles(topath)

		t.Assert(IsFile(testpath()+topath), true)
		t.AssertNE(CopyFile("", ""), nil)
	})
	// Content replacement.
	gtest.C(t, func(t *gtest.T) {
		src := TempDir(gtime.TimestampNanoStr())
		dst := TempDir(gtime.TimestampNanoStr())
		srcContent := "1"
		dstContent := "1"
		t.Assert(PutContents(src, srcContent), nil)
		t.Assert(PutContents(dst, dstContent), nil)
		t.Assert(GetContents(src), srcContent)
		t.Assert(GetContents(dst), dstContent)

		t.Assert(CopyFile(src, dst), nil)
		t.Assert(GetContents(src), srcContent)
		t.Assert(GetContents(dst), srcContent)
	})
}

func Test_CopyDir(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			dirPath1 = "/test-copy-dir1"
			dirPath2 = "/test-copy-dir2"
		)
		haveList := []string{
			"t1.txt",
			"t2.txt",
		}
		createDir(dirPath1)
		for _, v := range haveList {
			t.Assert(createTestFile(dirPath1+"/"+v, ""), nil)
		}
		defer delTestFiles(dirPath1)

		var (
			yfolder  = testpath() + dirPath1
			tofolder = testpath() + dirPath2
		)

		if IsDir(tofolder) {
			t.Assert(Remove(tofolder), nil)
			t.Assert(Remove(""), nil)
		}

		t.Assert(CopyDir(yfolder, tofolder), nil)
		defer delTestFiles(tofolder)

		t.Assert(IsDir(yfolder), true)

		for _, v := range haveList {
			t.Assert(IsFile(yfolder+"/"+v), true)
		}

		t.Assert(IsDir(tofolder), true)

		for _, v := range haveList {
			t.Assert(IsFile(tofolder+"/"+v), true)
		}

		t.Assert(Remove(tofolder), nil)
		t.Assert(Remove(""), nil)
	})
	// Content replacement.
	gtest.C(t, func(t *gtest.T) {
		src := TempDir(gtime.TimestampNanoStr(), gtime.TimestampNanoStr())
		dst := TempDir(gtime.TimestampNanoStr(), gtime.TimestampNanoStr())
		defer func() {
			Remove(src)
			Remove(dst)
		}()
		srcContent := "1"
		dstContent := "1"
		t.Assert(PutContents(src, srcContent), nil)
		t.Assert(PutContents(dst, dstContent), nil)
		t.Assert(GetContents(src), srcContent)
		t.Assert(GetContents(dst), dstContent)

		err := CopyDir(Dir(src), Dir(dst))
		t.Assert(err, nil)
		t.Assert(GetContents(src), srcContent)
		t.Assert(GetContents(dst), srcContent)
	})
}
