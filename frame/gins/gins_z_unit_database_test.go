// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gins_test

import (
	"github.com/basicfu/gf/debug/gdebug"
	"github.com/basicfu/gf/os/gtime"
	"testing"
	"time"

	"github.com/basicfu/gf/os/gfile"
	"github.com/basicfu/gf/test/gtest"
)

func Test_Database(t *testing.T) {
	databaseContent := gfile.GetContents(
		gdebug.TestDataPath("database", "config.toml"),
	)
	gtest.C(t, func(t *gtest.T) {
		var err error
		dirPath := gfile.TempDir(gtime.TimestampNanoStr())
		err = gfile.Mkdir(dirPath)
		t.Assert(err, nil)
		defer gfile.Remove(dirPath)

		name := "config.toml"
		err = gfile.PutContents(gfile.Join(dirPath, name), databaseContent)
		t.Assert(err, nil)

		err = Config().AddPath(dirPath)
		t.Assert(err, nil)

		defer Config().Clear()

		// for gfsnotify callbacks to refresh cache of config file
		time.Sleep(500 * time.Millisecond)

		//fmt.Println("gins Test_Database", Config().Get("test"))

		dbDefault := Database()
		dbTest := Database("test")
		t.AssertNE(dbDefault, nil)
		t.AssertNE(dbTest, nil)

		t.Assert(dbDefault.PingMaster(), nil)
		t.Assert(dbDefault.PingSlave(), nil)
		t.Assert(dbTest.PingMaster(), nil)
		t.Assert(dbTest.PingSlave(), nil)
	})
}
