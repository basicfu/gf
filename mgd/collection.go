package mgd

import (
	"context"
	"reflect"

	"github.com/basicfu/gf/g"
	"github.com/basicfu/gf/mgd/builder"
	"github.com/basicfu/gf/mgd/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection[T any | g.Map] struct {
	coll  *mongo.Collection
	model T
}

func (c *Collection[T]) FindOne(ctx context.Context, filter any) T {
	return c.FindOneByExample(ctx, Example{Filter: filter})
}
func (c *Collection[T]) FindOneByExample(ctx context.Context, example Example) T {
	opt := findOneOptions(example)
	m := *new(T)
	result := c.coll.FindOne(ctx, example.Filter, &opt)
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
func (c *Collection[T]) FindById(ctx context.Context, id any) T {
	return c.FindOneByExample(ctx, Example{Filter: g.Map{field.ID: Id(id)}})
}
func (c *Collection[T]) FindByIds(ctx context.Context, ids any) []T {
	return c.FindByExample(ctx, Example{Filter: g.Map{field.ID: g.Map{"$in": Ids(ids)}}})
}

func (c *Collection[T]) FindByExample(ctx context.Context, example Example) []T {
	opt := findOptions(example)
	m := make([]T, 0)
	cur, err := c.coll.Find(ctx, example.Filter, &opt)
	if err != nil {
		panic(err.Error())
	}
	err = cur.All(ctx, &m)
	if err != nil {
		panic(err.Error())
	}
	return m
}
func (c *Collection[T]) Find(ctx context.Context, filter any) []T {
	return c.FindByExample(ctx, Example{Filter: filter})
}
func (c *Collection[T]) FindAll(ctx context.Context) []T {
	return c.FindByExample(ctx, Example{Filter: g.Map{}})
}

func (c *Collection[T]) FindPageByExample(ctx context.Context, example Example) PageList[T] {
	f := findOptions(example)
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
	filter := example.Filter
	total := c.Count(ctx, filter)
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
func (c *Collection[T]) FindPage(ctx context.Context, filter any) PageList[T] {
	return c.FindPageByExample(ctx, Example{Filter: filter})
}

func (c *Collection[T]) FindOneAndUpdate(ctx context.Context, opt UpdateOptions, r interface{}) bool {
	op := options.FindOneAndUpdateOptions{}
	op.SetUpsert(opt.Upsert)
	if !opt.ReturnOldDocument { //默认返回更新后的文档
		op.SetReturnDocument(options.After)
	}
	if opt.Select != nil || opt.Exclude != nil {
		var projection bson.D
		for _, v := range opt.Select {
			projection = append(projection, bson.E{Key: v, Value: 1})
		}
		for _, v := range opt.Exclude {
			projection = append(projection, bson.E{Key: v, Value: 0})
		}
		op.Projection = projection
	}
	update, _ := updateOptions(opt)
	result := c.coll.FindOneAndUpdate(ctx, opt.Filter, update, &op)
	if result.Err() != nil {
		reflect.ValueOf(r).Elem().FieldByName("Nil").SetBool(true)
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

// ===============================
func (c *Collection[T]) Count(ctx context.Context, filter any) int64 {
	count, err := c.coll.CountDocuments(ctx, filter)
	if err != nil {
		panic(err.Error())
	}
	return count
}

// ===============================
func (c *Collection[T]) Insert(ctx context.Context, model any) interface{} {
	Create(model) //model非&时无法写入时间
	res, err := c.coll.InsertOne(ctx, model)
	if err != nil {
		panic(err)
	}
	return res.InsertedID
}

// 批量添加，不能超过isMaster.maxWriteBatchSize默认值10w条
func (c *Collection[T]) InsertMany(ctx context.Context, documents []any) []interface{} {
	var doc []any
	for _, v := range documents {
		Create(v)
		doc = append(doc, v)
	}
	i := options.InsertManyOptions{}
	res, err := c.coll.InsertMany(ctx, doc, &i)
	if err != nil {
		panic(err)
	}
	return res.InsertedIDs
}

// ===============================
func (c *Collection[T]) UpdateOne(ctx context.Context, opt UpdateOptions) mongo.UpdateResult {
	update, op := updateOptions(opt)
	updateResult, err := c.coll.UpdateOne(ctx, opt.Filter, update, &op)
	if err != nil {
		panic(err)
	}
	return *updateResult
}
func (c *Collection[T]) UpdateMany(ctx context.Context, opt UpdateOptions) mongo.UpdateResult {
	update, op := updateOptions(opt)
	updateResult, err := c.coll.UpdateMany(ctx, opt.Filter, update, &op)
	if err != nil {
		panic(err)
	}
	return *updateResult
}
func updateOptions(opt UpdateOptions) (bson.M, options.UpdateOptions) {
	op := options.UpdateOptions{}
	op.SetUpsert(opt.Upsert)
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
	return update, op
}

// ===============================
func (c *Collection[T]) Delete(ctx context.Context, filter any) int64 {
	res, err := c.coll.DeleteMany(ctx, filter)
	if err != nil {
		panic(err)
	}
	return res.DeletedCount
}
func (c *Collection[T]) DeleteOne(ctx context.Context, filter any) int64 {
	res, err := c.coll.DeleteOne(ctx, filter)
	if err != nil {
		panic(err)
	}
	return res.DeletedCount
}
func (c *Collection[T]) DeleteById(ctx context.Context, id any) int64 {
	return c.DeleteByIds(ctx, []any{id})
}
func (c *Collection[T]) DeleteByIds(ctx context.Context, ids []any) int64 {
	var res *mongo.DeleteResult
	var err error
	if len(ids) == 1 {
		res, err = c.coll.DeleteOne(ctx, bson.M{field.ID: Id(ids[0])})
	} else {
		res, err = c.coll.DeleteMany(ctx, bson.M{field.ID: bson.M{"$in": ids}})
	}
	if err != nil {
		panic(err)
	}
	return res.DeletedCount
}

// ===============================
func (c *Collection[T]) SimpleAggregate(ctx context.Context, results interface{}, stages ...interface{}) error {
	cur, err := c.simpleAggregateCursor(ctx, stages...)
	if err != nil {
		return err
	}
	return cur.All(ctx, results)
}
func (c *Collection[T]) simpleAggregateCursor(ctx context.Context, stages ...interface{}) (*mongo.Cursor, error) {
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
