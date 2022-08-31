package mgd

import (
	"github.com/basicfu/gf/errors/gerror"
	"go.mongodb.org/mongo-driver/mongo"
)

//TODO 可以考虑在事物开始时使用defer
func Transaction(callback func(ctx mongo.SessionContext)) {
	session, e := client.StartSession()
	if e != nil {
		panic(e)
	}
	ctx := ctx()
	defer session.EndSession(ctx)
	_, e = session.WithTransaction(ctx, func(context mongo.SessionContext) (d interface{}, err error) {
		defer func() {
			if errRec := recover(); errRec != nil {
				//这里没办法做成通用方法，除非把exception.error部分提取到gerror
				switch errRec.(type) {
				case gerror.Error:
					err = errRec.(gerror.Error)
				case error:
					err = errRec.(error)
				}
			}
		}()
		callback(context)
		return nil, nil
	})
	if e != nil {
		panic(e)
	}
}
