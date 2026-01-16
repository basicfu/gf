package g

import (
	"fmt"

	"github.com/basicfu/gf/decimal"
)

// 全局错误无法捕捉线程内的，使用线程时务必处理错误避免挂掉整个服务
var Go = func(handler func(), catch ...func(err error)) {
	go Try(handler, catch...)
}
var Func = func(handler func()) {
	handler()
}
var Decimal = decimal.New

type Map = map[string]interface{}

func Try(try func(), catch ...func(err error)) {
	defer func() {
		if err := recover(); err != nil && len(catch) > 0 {
			if err, ok := err.(error); ok {
				catch[0](err)
			} else {
				catch[0](fmt.Errorf(`%v`, err))
			}
		}
	}()
	try()
}

func Recover(catch ...func(err error)) {
	if err := recover(); err != nil {
		if len(catch) == 0 {
			return
		}
		if v, ok := err.(error); ok {
			catch[0](v)
		} else {
			catch[0](fmt.Errorf(`%v`, err))
		}
	}
}
