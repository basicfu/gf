// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gfile_test

import (
	"github.com/basicfu/gf/container/garray"
	"github.com/basicfu/gf/debug/gdebug"
	"testing"

	"github.com/basicfu/gf/test/gtest"
)

func Test_ScanDir(t *testing.T) {
	teatPath := gdebug.TestDataPath()
	gtest.C(t, func(t *gtest.T) {
		files, err := ScanDir(teatPath, "*", false)
		t.Assert(err, nil)
		t.AssertIN(teatPath+Separator+"dir1", files)
		t.AssertIN(teatPath+Separator+"dir2", files)
		t.AssertNE(teatPath+Separator+"dir1"+Separator+"file1", files)
	})
	gtest.C(t, func(t *gtest.T) {
		files, err := ScanDir(teatPath, "*", true)
		t.Assert(err, nil)
		t.AssertIN(teatPath+Separator+"dir1", files)
		t.AssertIN(teatPath+Separator+"dir2", files)
		t.AssertIN(teatPath+Separator+"dir1"+Separator+"file1", files)
		t.AssertIN(teatPath+Separator+"dir2"+Separator+"file2", files)
	})
}

func Test_ScanDirFunc(t *testing.T) {
	teatPath := gdebug.TestDataPath()
	gtest.C(t, func(t *gtest.T) {
		files, err := ScanDirFunc(teatPath, "*", true, func(path string) string {
			if Name(path) != "file1" {
				return ""
			}
			return path
		})
		t.Assert(err, nil)
		t.Assert(len(files), 1)
		t.Assert(Name(files[0]), "file1")
	})
}

func Test_ScanDirFile(t *testing.T) {
	teatPath := gdebug.TestDataPath()
	gtest.C(t, func(t *gtest.T) {
		files, err := ScanDirFile(teatPath, "*", false)
		t.Assert(err, nil)
		t.Assert(len(files), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		files, err := ScanDirFile(teatPath, "*", true)
		t.Assert(err, nil)
		t.AssertNI(teatPath+Separator+"dir1", files)
		t.AssertNI(teatPath+Separator+"dir2", files)
		t.AssertIN(teatPath+Separator+"dir1"+Separator+"file1", files)
		t.AssertIN(teatPath+Separator+"dir2"+Separator+"file2", files)
	})
}

func Test_ScanDirFileFunc(t *testing.T) {
	teatPath := gdebug.TestDataPath()
	gtest.C(t, func(t *gtest.T) {
		array := garray.New()
		files, err := ScanDirFileFunc(teatPath, "*", false, func(path string) string {
			array.Append(1)
			return path
		})
		t.Assert(err, nil)
		t.Assert(len(files), 0)
		t.Assert(array.Len(), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.New()
		files, err := ScanDirFileFunc(teatPath, "*", true, func(path string) string {
			array.Append(1)
			if Basename(path) == "file1" {
				return path
			}
			return ""
		})
		t.Assert(err, nil)
		t.Assert(len(files), 1)
		t.Assert(array.Len(), 3)
	})
}
