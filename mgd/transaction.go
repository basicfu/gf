package mgd

import (
	"context"
)

// 每个事物并发时都会重复执行，同时进行中的事物，会在一个事物完成后另一个事物重新执行一遍，业务时需要做好处理
func Transaction(ctx context.Context, callback func(ctx context.Context)) {
	session, e := client.StartSession()
	if e != nil {
		panic(e)
	}
	defer func() {
		session.EndSession(ctx) //事务内错误页也会执行defer中止事务
	}()
	_, e = session.WithTransaction(ctx, func(context context.Context) (d interface{}, err error) {
		callback(context)
		return nil, nil
	})
}
