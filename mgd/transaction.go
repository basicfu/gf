package mgd

import (
	"go.mongodb.org/mongo-driver/mongo"
)

//func Transaction(callback func(ctx mongo.SessionContext)) {
//	session, err := client.StartSession()
//	if err != nil {
//		panic(err)
//	}
//	defer session.EndSession(Ctx())
//	result, err := session.WithTransaction(Ctx(), func(sessCtx mongo.SessionContext) (i interface{}, e error) {
//		defer func() {
//			if err := recover(); err != nil {
//				i = err //使用interface抛出自定义错
//				e = errors.New("")
//			}
//		}()
//		callback(sessCtx)
//		return i, e
//	}, options.Transaction().
//		SetReadConcern(readconcern.Snapshot()).
//		SetWriteConcern(writeconcern.New(writeconcern.WMajority())))
//	if result != nil {
//		panic(result)
//	} else if err != nil {
//		panic(err)
//	}
//}

func Transaction(callback func(ctx mongo.SessionContext)) {
	session, err := client.StartSession()
	ctx := ctx()
	defer session.EndSession(ctx)
	session.StartTransaction()
	mongo.WithSession(ctx, session, func(context mongo.SessionContext) error {
		defer func() {
			if err := recover(); err != nil {
				//err = errors.New("")
				//TODO
				session.AbortTransaction(ctx) //有异常，事物终止
				panic(err)
			}
		}()
		callback(context)
		session.CommitTransaction(ctx) //提交事物
		return err
	})

}
