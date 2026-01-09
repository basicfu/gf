package g

import (
	"fmt"

	"github.com/basicfu/gf/decimal"
)

// 全局错误无法捕捉线程内的，使用线程时务必处理错误避免挂掉整个服务
var Go = func(handler func(), catch ...func(err error)) {
	go Try(handler, catch...)
}

var Decimal = decimal.New

type Map = map[string]interface{}

func Try(try func(), catch ...func(err error)) {
	defer func() {
		if exception := recover(); exception != nil && len(catch) > 0 {
			if err, ok := exception.(error); ok {
				catch[0](err)
			} else {
				catch[0](fmt.Errorf(`%v`, exception))
			}
		}
	}()
	try()
}
