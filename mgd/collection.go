package mgd

import (
	"context"
	"github.com/basicfu/gf/g"
	"github.com/basicfu/gf/mgd/builder"
	"github.com/basicfu/gf/mgd/field"
	"github.com/basicfu/gf/util/gconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

type Collection[T any | g.Map] struct {
	coll  *mongo.Collection
	model T
}

//-----findOne------
func (c *Collection[T]) FindOneByExample(example Example) T {
	opt := findOneOptions(example)
	m := *new(T)
	ctx := example.Context
	if ctx == nil {
		ctx = buildCtx()
	}
	result := c.coll.FindOne(ctx, toFilter(example.Filter), &opt)
	if result.Err() != nil {
		if mongo.ErrNoDocuments.Error() == result.Err().Error() {
			reflect.ValueOf(&m).Elem().FieldByName("Nil").SetBool(true) //标识对象业务为空
			return m
		} else {
			panic(result.Err())
		}
	}
	err := result.Decode(&m)
	if err != nil {
		panic(err.Error())
	}
	return m
}
func (c *Collection[T]) FindOne(filter any, ctxArray ...context.Context) T {
	return c.FindOneByExample(Example{Context: buildCtx(ctxArray...), Filter: filter})
}

func (c *Collection[T]) FindById(id any, ctxArray ...context.Context) T {
	return c.FindOneByExample(Example{Context: buildCtx(ctxArray...), Filter: g.Map{field.ID: Id(id)}})
}

//---------------
//func (c *Collection) FindByIds(ids interface{}, result interface{}) {
//	c.Find(FindOptions{
//		Filter: bson.M{field.ID: bson.M{"$in": Ids(ids)}},
//	}, result)
//}

func (c *Collection[T]) Find(opt FindOptions, result interface{}) {
	f := findOptions(&opt)
	cur, err := c.coll.Find(buildCtx(), opt.Filter, &f)
	if err != nil {
		panic(err.Error()) //TODO 超时错误处理，应该在全局错误处详细捕捉
	}
	err = cur.All(buildCtx(), result) //可考虑如果不传类型使用coll创建时默认的
	if err != nil {
		panic(err.Error())
	}
}
func (c *Collection[T]) FindAll(result interface{}) {
	c.Find(FindOptions{}, result)
}
func (c *Collection[T]) FindResult(opt FindOptions) interface{} {
	c.Find(opt, &opt.Result)
	return opt.Result
}
func (c *Collection[T]) FindPageResult(opt FindOptions) PageList {
	f := findOptions(&opt)
	page := Page{}
	if opt.Page.PageSize != 0 && opt.Page.PageNum != 0 {
		total := c.Count(opt.Filter)
		if total == 0 {
			return PageList{List: []string{}}
		}
		page.PageSize = opt.Page.PageSize
		page.PageNum = opt.Page.PageNum
		page.Total = total
		maxPage := total / page.PageSize
		if total%page.PageSize != 0 {
			maxPage = maxPage + 1
		}
		if page.PageNum > maxPage {
			page.PageNum = maxPage
		}
		skip := (page.PageNum - 1) * page.PageSize
		f.Skip = &skip
		f.Limit = &page.PageSize
	}
	r := opt.Result
	if opt.Result == nil {
		t := reflect.SliceOf(reflect.TypeOf(c.model))
		r = reflect.MakeSlice(t, 0, 0).Interface()
	}
	cur, err := c.coll.Find(buildCtx(), opt.Filter, &f)
	if err != nil {
		panic(err)
	}
	err = cur.All(buildCtx(), &r)
	if err != nil {
		panic(err)
	}
	return PageList{
		Page: page,
		List: r,
	}
}
func (c *Collection[T]) Count(filter interface{}, opts ...*options.CountOptions) int64 {
	if filter == nil {
		filter = bson.M{}
	}
	count, err := c.coll.CountDocuments(buildCtx(), filter, opts...)
	if err != nil {
		panic(err.Error())
	}
	return count
}
func (c *Collection[T]) CountCtx(filter interface{}, ctx context.Context) int64 {
	if filter == nil {
		filter = bson.M{}
	}
	count, err := c.coll.CountDocuments(ctx, filter)
	if err != nil {
		panic(err.Error())
	}
	return count
}
func (c *Collection[T]) Insert(model interface{}, ctxArray ...context.Context) interface{} {
	Create(model)
	useCtx := buildCtx()
	if len(ctxArray) != 0 {
		useCtx = ctxArray[0] //事物
	}
	res, err := c.coll.InsertOne(useCtx, model)
	if err != nil {
		panic(err)
	}
	return res.InsertedID
}

//批量添加，不能超过isMaster.maxWriteBatchSize默认值10w条
func (c *Collection[T]) InsertMany(opt InsertOptions) []interface{} {
	doc := []interface{}{}
	for _, v := range gconv.SliceAny(opt.Document) {
		Create(v)
		doc = append(doc, v)
	}
	i := options.InsertManyOptions{}
	if opt.Context == nil {
		opt.Context = buildCtx()
	}
	res, err := c.coll.InsertMany(opt.Context, doc, &i)
	if err != nil {
		panic(err)
	}
	return res.InsertedIDs
}
func (c *Collection[T]) FindOneAndUpdate(opt UpdateOptions, r interface{}) bool {
	updateOptions := options.FindOneAndUpdateOptions{}
	updateOptions.SetUpsert(opt.Upset)
	//if opt.ReturnNewDocument {//默认为true
	updateOptions.SetReturnDocument(options.After)
	//}
	update := bson.M{}
	if opt.Set != nil {
		if hook, ok := opt.Set.(UpdateHook); ok {
			hook.Update(nil)
		}
		update["$set"] = opt.Set
	}
	if opt.Inc != nil {
		update["$inc"] = opt.Inc
	}
	result := c.coll.FindOneAndUpdate(opt.Context, opt.Filter, update, &updateOptions)
	if result.Err() != nil {
		if mongo.ErrNoDocuments.Error() == result.Err().Error() {
			return false
		} else {
			panic(result.Err())
		}
	}
	err := result.Decode(r)
	if err != nil {
		panic(err)
	}
	return true
}
func (c *Collection[T]) UpdateOne(opt UpdateOptions) mongo.UpdateResult {
	updateOptions := options.UpdateOptions{}
	updateOptions.SetUpsert(opt.Upset)
	update := bson.M{}
	if opt.Set != nil {
		if hook, ok := opt.Set.(UpdateHook); ok { //如果使用Update类型自动更新时间
			hook.Update(nil)
		}
		update["$set"] = opt.Set
	}
	if opt.Inc != nil {
		update["$inc"] = opt.Inc
	}
	if opt.UnSet != nil {
		update["$unset"] = opt.UnSet
	}
	if opt.AddToSet != nil {
		update["$addToSet"] = opt.AddToSet
	}
	if opt.Push != nil {
		update["$push"] = opt.Push
	}
	updateResult, err := c.coll.UpdateOne(opt.Context, opt.Filter, update, &updateOptions)
	if err != nil {
		panic(err)
	}
	return *updateResult
}
func (c *Collection[T]) UpdateMany(opt UpdateOptions) mongo.UpdateResult {
	updateOptions := options.UpdateOptions{}
	updateOptions.SetUpsert(opt.Upset)
	update := bson.M{}
	if opt.Set != nil {
		if hook, ok := opt.Set.(UpdateHook); ok { //如果使用Update类型自动更新时间
			hook.Update(nil)
		}
		update["$set"] = opt.Set
	}
	if opt.Inc != nil {
		update["$inc"] = opt.Inc
	}
	if opt.UnSet != nil {
		update["$unset"] = opt.UnSet
	}
	if opt.AddToSet != nil {
		update["$addToSet"] = opt.AddToSet
	}
	if opt.Push != nil {
		update["$push"] = opt.Push
	}
	updateResult, err := c.coll.UpdateMany(opt.Context, opt.Filter, update, &updateOptions)
	if err != nil {
		panic(err)
	}
	return *updateResult
}
func (c *Collection[T]) UpdateById(model Model, opts ...*options.UpdateOptions) {
	Update(model)
	_, err := c.coll.UpdateOne(buildCtx(), bson.M{field.ID: model.GetId()}, bson.M{"$set": model}, opts...)
	if err != nil {
		panic(err)
	}
}

func (c *Collection[T]) Delete(opt DeleteOptions) int64 {
	var res *mongo.DeleteResult
	var err error
	if opt.Filter == nil {
		opt.Filter = bson.M{}
	}
	if opt.Context == nil {
		opt.Context = buildCtx()
	}
	res, err = c.coll.DeleteMany(opt.Context, opt.Filter)
	if err != nil {
		panic(err)
	}
	return res.DeletedCount
}
func (c *Collection[T]) DeleteByIds(ids []string) int64 {
	var res *mongo.DeleteResult
	var err error
	if len(ids) == 1 {
		res, err = c.coll.DeleteOne(buildCtx(), bson.M{field.ID: Id(ids[0])})
	} else {
		res, err = c.coll.DeleteMany(buildCtx(), bson.M{field.ID: bson.M{"$in": ids}})
	}
	if err != nil {
		panic(err)
	}
	return res.DeletedCount
}

func (c *Collection[T]) SimpleAggregateFirst(result interface{}, stages ...interface{}) (bool, error) {
	cur, err := c.SimpleAggregateCursor(stages...)
	if err != nil {
		return false, err
	}
	if cur.Next(buildCtx()) {
		return true, cur.Decode(result)
	}
	return false, nil
}
func (c *Collection[T]) SimpleAggregate(results interface{}, stages ...interface{}) error {
	cur, err := c.SimpleAggregateCursor(stages...)
	if err != nil {
		return err
	}
	return cur.All(buildCtx(), results)
}
func (c *Collection[T]) SimpleAggregateCursor(stages ...interface{}) (*mongo.Cursor, error) {
	pipeline := bson.A{}
	for _, stage := range stages {
		if operator, ok := stage.(builder.Operator); ok {
			pipeline = append(pipeline, builder.S(operator))
		} else {
			pipeline = append(pipeline, stage)
		}
	}
	return c.coll.Aggregate(buildCtx(), pipeline, nil)
}
