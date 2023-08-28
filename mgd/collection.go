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

// -----findOne------
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

// filter只允许g.map和struct结构，但是目前没法限制只传入这两个类型
func (c *Collection[T]) FindOne(filter any, ctxArray ...context.Context) T {
	return c.FindOneByExample(Example{Context: buildCtx(ctxArray...), Filter: filter})
}

func (c *Collection[T]) FindById(id any, ctxArray ...context.Context) T {
	return c.FindOneByExample(Example{Context: buildCtx(ctxArray...), Filter: g.Map{field.ID: Id(id)}})
}
func (c *Collection[T]) FindByIds(ids any, ctxArray ...context.Context) []T {
	return c.FindByExample(Example{Context: buildCtx(ctxArray...), Filter: g.Map{field.ID: g.Map{"$in": Ids(ids)}}})
}

// -----find-----
func (c *Collection[T]) FindByExample(example Example) []T {
	opt := findOptions(example)
	m := make([]T, 0)
	ctx := example.Context
	if ctx == nil {
		ctx = buildCtx()
	}
	cur, err := c.coll.Find(ctx, toFilter(example.Filter), &opt)
	if err != nil {
		panic(err.Error())
	}
	err = cur.All(ctx, &m)
	if err != nil {
		panic(err.Error())
	}
	return m
}
func (c *Collection[T]) FindByExampleTest(example Example) []g.Map {
	opt := findOptions(example)
	m := make([]g.Map, 0)
	ctx := example.Context
	if ctx == nil {
		ctx = buildCtx()
	}
	cur, err := c.coll.Find(ctx, toFilter(example.Filter), &opt)
	if err != nil {
		panic(err.Error())
	}
	err = cur.All(ctx, &m)
	if err != nil {
		panic(err.Error())
	}
	return m
}
func (c *Collection[T]) Find(filter any, ctxArray ...context.Context) []T {
	return c.FindByExample(Example{Context: buildCtx(ctxArray...), Filter: filter})
}
func (c *Collection[T]) FindAll(ctxArray ...context.Context) []T {
	return c.FindByExample(Example{Context: buildCtx(ctxArray...), Filter: g.Map{}})
}

// -----findPage-----
func (c *Collection[T]) FindPage(filter any, ctxArray ...context.Context) PageList[T] {
	return c.FindPageByExample(Example{Context: buildCtx(ctxArray...), Filter: filter})
}
func (c *Collection[T]) FindPageByExample(example Example) PageList[T] {
	f := findOptions(example)
	ctx := example.Context
	if ctx == nil {
		ctx = buildCtx()
	}
	page := Page{}
	list := make([]T, 0)
	pageNum := example.Page.PageNum
	pageSize := example.Page.PageSize
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	filter := toFilter(example.Filter)
	total := c.Count(filter, ctx)
	if total == 0 {
		return PageList[T]{List: list}
	}
	page.PageSize = pageSize
	page.PageNum = pageNum
	page.Total = total
	maxPage := total / page.PageSize
	if total%page.PageSize != 0 {
		maxPage = maxPage + 1
	}
	//if page.PageNum > maxPage {
	//	page.PageNum = maxPage
	//}
	skip := (page.PageNum - 1) * page.PageSize
	f.Skip = &skip
	f.Limit = &page.PageSize
	cur, err := c.coll.Find(ctx, filter, &f)
	if err != nil {
		panic(err.Error())
	}
	err = cur.All(ctx, &list)
	if err != nil {
		panic(err.Error())
	}
	return PageList[T]{
		Page: page,
		List: list,
	}
}

// -----count------
func (c *Collection[T]) Count(filter any, ctxArray ...context.Context) int64 {
	count, err := c.coll.CountDocuments(buildCtx(ctxArray...), toFilter(filter))
	if err != nil {
		panic(err.Error())
	}
	return count
}

// -----insert------
func (c *Collection[T]) Insert(model interface{}, ctxArray ...context.Context) interface{} {
	Create(model) //model非&时无法写入时间
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

// 批量添加，不能超过isMaster.maxWriteBatchSize默认值10w条
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
	updateOptions.SetUpsert(opt.Upsert)
	if !opt.ReturnOldDocument { //默认返回更新后的文档
		updateOptions.SetReturnDocument(options.After)
	}
	if opt.Select != nil || opt.Exclude != nil {
		var projection bson.D
		for _, v := range opt.Select {
			projection = append(projection, bson.E{Key: v, Value: 1})
		}
		for _, v := range opt.Exclude {
			projection = append(projection, bson.E{Key: v, Value: 0})
		}
		updateOptions.Projection = projection
	}
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
	updateOptions.SetUpsert(opt.Upsert)
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
	if opt.Pull != nil {
		update["$pull"] = opt.Pull
	}
	updateResult, err := c.coll.UpdateOne(opt.Context, opt.Filter, update, &updateOptions)
	if err != nil {
		panic(err)
	}
	return *updateResult
}
func (c *Collection[T]) UpdateMany(opt UpdateOptions) mongo.UpdateResult {
	updateOptions := options.UpdateOptions{}
	updateOptions.SetUpsert(opt.Upsert)
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
	if opt.Pull != nil {
		update["$pull"] = opt.Pull
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
	if opt.Filter == nil {
		opt.Filter = bson.M{}
	}
	if opt.Context == nil {
		opt.Context = buildCtx()
	}
	var res *mongo.DeleteResult
	var err error
	if opt.One {
		res, err = c.coll.DeleteOne(opt.Context, opt.Filter)
	} else {
		res, err = c.coll.DeleteMany(opt.Context, opt.Filter)
	}
	if err != nil {
		panic(err)
	}
	return res.DeletedCount
}
func (c *Collection[T]) DeleteOne(opt DeleteOptions) int64 {
	opt.One = true
	return c.Delete(opt)
}
func (c *Collection[T]) DeleteById(id any, ctxArray ...context.Context) int64 {
	var res *mongo.DeleteResult
	var err error
	res, err = c.coll.DeleteOne(buildCtx(ctxArray...), bson.M{field.ID: Id(id)})
	if err != nil {
		panic(err)
	}
	return res.DeletedCount
}
func (c *Collection[T]) DeleteByIds(ids []any, ctxArray ...context.Context) int64 {
	var res *mongo.DeleteResult
	var err error
	if len(ids) == 1 {
		res, err = c.coll.DeleteOne(buildCtx(ctxArray...), bson.M{field.ID: Id(ids[0])})
	} else {
		res, err = c.coll.DeleteMany(buildCtx(ctxArray...), bson.M{field.ID: bson.M{"$in": ids}})
	}
	if err != nil {
		panic(err)
	}
	return res.DeletedCount
}

func (c *Collection[T]) SimpleAggregateFirst(result interface{}, stages ...interface{}) (bool, error) {
	cur, err := c.SimpleAggregateCursor(buildCtx(), stages...)
	if err != nil {
		return false, err
	}
	if cur.Next(buildCtx()) {
		return true, cur.Decode(result)
	}
	return false, nil
}
func (c *Collection[T]) SimpleAggregate(results interface{}, stages ...interface{}) error {
	ctx := buildCtx()
	cur, err := c.SimpleAggregateCursor(ctx, stages...)
	if err != nil {
		return err
	}
	return cur.All(ctx, results)
}
func (c *Collection[T]) SimpleAggregateCtx(ctx context.Context, results interface{}, stages ...interface{}) error {
	cur, err := c.SimpleAggregateCursor(ctx, stages...)
	if err != nil {
		return err
	}
	return cur.All(ctx, results)
}
func (c *Collection[T]) SimpleAggregateCursor(ctx context.Context, stages ...interface{}) (*mongo.Cursor, error) {
	pipeline := bson.A{}
	for _, stage := range stages {
		if operator, ok := stage.(builder.Operator); ok {
			pipeline = append(pipeline, builder.S(operator))
		} else {
			pipeline = append(pipeline, stage)
		}
	}
	return c.coll.Aggregate(ctx, pipeline, nil)
}
