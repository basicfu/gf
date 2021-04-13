// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

package gjson_test

import (
	"fmt"
	"github.com/basicfu/gf/debug/gdebug"
)

func Example_loadJson() {
	jsonFilePath := gdebug.TestDataPath("json", "data1.json")
	j, _ := Load(jsonFilePath)
	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))
}

func Example_loadXml() {
	jsonFilePath := gdebug.TestDataPath("xml", "data1.xml")
	j, _ := Load(jsonFilePath)
	fmt.Println(j.Get("doc.name"))
	fmt.Println(j.Get("doc.score"))
}

func Example_loadContent() {
	jsonContent := `{"name":"john", "score":"100"}`
	j, _ := LoadContent(jsonContent)
	fmt.Println(j.Get("name"))
	fmt.Println(j.Get("score"))
	// Output:
	// john
	// 100
}
