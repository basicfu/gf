package mgd

import (
	"context"
	"errors"
	"fmt"

	"github.com/basicfu/gf/gerror"
	"go.mongodb.org/mongo-driver/mongo"
)

// 同时进行中的事物，会在一个事物完成后另一个事物重新执行一遍，业务时需要做好处理
// 每个事物并发时都会重复执行
func Transaction(ctx context.Context, callback func(ctx mongo.SessionContext)) {
	session, e := client.StartSession()
	if e != nil {
		panic(e)
	}
	defer session.EndSession(ctx)
	_, e = session.WithTransaction(ctx, func(context mongo.SessionContext) (d interface{}, err error) {
		defer func() {
			if errRec := recover(); errRec != nil {
				switch v := errRec.(type) {
				case gerror.Error:
					err = v
				case error:
					err = v
				case string:
					err = errors.New(v)
				default:
					err = fmt.Errorf("panic: %v", v)
				}
			}
		}()
		callback(context)
		return nil, nil
	})
	if e != nil {
		panic(gerror.Msg("Transaction failed").WithError(e).WithSkip(2))
	}
}
