package mgd

import (
	"context"
	"github.com/basicfu/gf/mgd/builder"
	"github.com/basicfu/gf/mgd/field"
	"github.com/basicfu/gf/util/gconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

type Collection struct {
	coll  *mongo.Collection
	model interface{}
}

func Id(id string) primitive.ObjectID {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return objectId
}
func PrepareId(id interface{}) interface{} {
	if idStr, ok := id.(string); ok {
		objectId, _ := primitive.ObjectIDFromHex(idStr)
		return objectId
	}
	return id
}
func PrepareIds(ids interface{}) interface{} {
	if idStr, ok := ids.([]string); ok {
		var newIds []primitive.ObjectID
		for _, v := range idStr {
			objectId, _ := primitive.ObjectIDFromHex(v)
			newIds = append(newIds, objectId)
		}

		return newIds
	}
	return ids
}

//TODO 想统一使用opt形式，return对象强转，true/false是否强制抛出错误
//FindById(id,result)
//FindByIds(ids,result)
//Find(filter,result)
//FindAll(result)
//FindByExample(example,result)
//FindOne(filter,result)
//FindOneByExample(example,result)
func (c *Collection) FindById(id interface{}, result interface{}) bool {
	return c.FindOneByExample(FindOptions{
		Filter: bson.M{field.ID: PrepareId(id)},
	}, result)
}
func (c *Collection) FindByIds(ids interface{}, result interface{}) {
	c.Find(FindOptions{
		Filter: bson.M{field.ID: bson.M{"$in": PrepareIds(ids)}},
	}, result)
}

//TODO
func (c *Collection) FindByIdResult(id interface{}) interface{} {
	var result interface{}
	c.FindOne(bson.M{field.ID: PrepareId(id)}, &result)
	return result
}
func (c *Collection) FindOne(filter interface{}, result interface{}, ctxArray ...context.Context) bool {
	useCtx := ctx()
	if len(ctxArray) != 0 {
		useCtx = ctxArray[0] //事物
	}
	return c.FindOneByExample(FindOptions{
		Context: useCtx,
		Filter:  filter,
	}, result)
}

//true为找到，false没找到//TODO 应该加加入example复杂查询类型，以及查询指定字段
func (c *Collection) FindOneResult(filter interface{}) (interface{}, bool) {
	result := reflect.New(reflect.TypeOf(c.model)).Interface()
	flag := c.FindOneByExample(FindOptions{
		Filter: filter,
	}, result)
	if !flag {
		result = bson.M{}
	}
	return result, flag
}
func (c *Collection) FindOneByExampleResult(opt FindOptions) (interface{}, bool) {
	flag := c.FindOneByExample(opt, opt.Result)
	if !flag {
		opt.Result = bson.M{}
	}
	return opt.Result, flag
}
func (c *Collection) FindOneByExample(opt FindOptions, result interface{}) bool {
	f := findOneOptions(&opt)
	useCtx := opt.Context
	if useCtx == nil {
		useCtx = ctx()
	}
	one := c.coll.FindOne(useCtx, opt.Filter, &f)
	if one.Err() != nil {
		if mongo.ErrNoDocuments.Error() == one.Err().Error() {
			if opt.NoFoundError {
				panic(one.Err())
			}
			return false
		} else {
			panic(one.Err())
		}
	}
	err := one.Decode(result)
	if err != nil {
		panic(err.Error())
	}
	return true
}
func (c *Collection) Find(opt FindOptions, result interface{}) {
	f := findOptions(&opt)
	cur, err := c.coll.Find(ctx(), opt.Filter, &f)
	if err != nil {
		panic(err.Error())
	}
	err = cur.All(ctx(), result) //可考虑如果不传类型使用coll创建时默认的
	if err != nil {
		panic(err.Error())
	}
}
func (c *Collection) FindAll(result interface{}) {
	c.Find(FindOptions{}, result)
}
func (c *Collection) FindResult(opt FindOptions) interface{} {
	c.Find(opt, &opt.Result)
	return opt.Result
}
func (c *Collection) FindPageResult(opt FindOptions) PageList {
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
	cur, err := c.coll.Find(ctx(), opt.Filter, &f)
	if err != nil {
		panic(err)
	}
	err = cur.All(ctx(), &r)
	if err != nil {
		panic(err)
	}
	return PageList{
		Page: page,
		List: r,
	}
}
func (c *Collection) Count(filter interface{}, opts ...*options.CountOptions) int64 {
	if filter == nil {
		filter = bson.M{}
	}
	count, err := c.coll.CountDocuments(ctx(), filter, opts...)
	if err != nil {
		panic(err.Error())
	}
	return count
}
func (c *Collection) Insert(model interface{}, ctxArray ...context.Context) interface{} {
	Create(model)
	useCtx := ctx()
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
func (c *Collection) InsertMany(opt InsertOptions) []interface{} {
	doc := []interface{}{}
	for _, v := range gconv.SliceAny(opt.Document) {
		Create(v)
		doc = append(doc, v)
	}
	i := options.InsertManyOptions{}
	if opt.Context == nil {
		opt.Context = ctx()
	}
	res, err := c.coll.InsertMany(opt.Context, doc, &i)
	if err != nil {
		panic(err)
	}
	return res.InsertedIDs
}
func (c *Collection) FindOneAndUpdate(opt UpdateOptions, r interface{}) bool {
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
func (c *Collection) UpdateOne(opt UpdateOptions) mongo.UpdateResult {
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
func (c *Collection) UpdateMany(opt UpdateOptions) mongo.UpdateResult {
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
func (c *Collection) UpdateById(model Model, opts ...*options.UpdateOptions) {
	Update(model)
	_, err := c.coll.UpdateOne(ctx(), bson.M{field.ID: model.GetId()}, bson.M{"$set": model}, opts...)
	if err != nil {
		panic(err)
	}
}

func (c *Collection) Delete(opt DeleteOptions) int64 {
	var res *mongo.DeleteResult
	var err error
	if opt.Filter == nil {
		opt.Filter = bson.M{}
	}
	if opt.Context == nil {
		opt.Context = ctx()
	}
	res, err = c.coll.DeleteMany(opt.Context, opt.Filter)
	if err != nil {
		panic(err)
	}
	return res.DeletedCount
}
func (c *Collection) DeleteByIds(ids []string) int64 {
	var res *mongo.DeleteResult
	var err error
	if len(ids) == 1 {
		res, err = c.coll.DeleteOne(ctx(), bson.M{field.ID: PrepareId(ids[0])})
	} else {
		res, err = c.coll.DeleteMany(ctx(), bson.M{field.ID: bson.M{"$in": ids}})
	}
	if err != nil {
		panic(err)
	}
	return res.DeletedCount
}

func (c *Collection) SimpleAggregateFirst(result interface{}, stages ...interface{}) (bool, error) {
	cur, err := c.SimpleAggregateCursor(stages...)
	if err != nil {
		return false, err
	}
	if cur.Next(ctx()) {
		return true, cur.Decode(result)
	}
	return false, nil
}
func (c *Collection) SimpleAggregate(results interface{}, stages ...interface{}) error {
	cur, err := c.SimpleAggregateCursor(stages...)
	if err != nil {
		return err
	}
	return cur.All(ctx(), results)
}
func (c *Collection) SimpleAggregateCursor(stages ...interface{}) (*mongo.Cursor, error) {
	pipeline := bson.A{}
	for _, stage := range stages {
		if operator, ok := stage.(builder.Operator); ok {
			pipeline = append(pipeline, builder.S(operator))
		} else {
			pipeline = append(pipeline, stage)
		}
	}
	return c.coll.Aggregate(ctx(), pipeline, nil)
}
func findOptions(opt *FindOptions) options.FindOptions {
	f := options.FindOptions{}
	if opt.Filter == nil {
		opt.Filter = bson.M{}
	}
	if opt.Limit != 0 {
		f.SetLimit(opt.Limit)
	}
	if opt.Asc != nil || opt.Desc != nil {
		var sort bson.D
		for _, v := range opt.Asc {
			sort = append(sort, bson.E{Key: v, Value: 1})
		}
		for _, v := range opt.Desc {
			sort = append(sort, bson.E{Key: v, Value: -1})
		}
		f.Sort = sort
	}
	if opt.Select != nil || opt.Exclude != nil {
		var projection bson.D
		for _, v := range opt.Select {
			projection = append(projection, bson.E{Key: v, Value: 1})
		}
		for _, v := range opt.Desc {
			projection = append(projection, bson.E{Key: v, Value: 0})
		}
		f.Projection = projection
	}
	if opt.BatchSize != 0 {
		f.SetBatchSize(opt.BatchSize)
	}
	return f
}
func findOneOptions(opt *FindOptions) options.FindOneOptions {
	f := options.FindOneOptions{}
	if opt.Filter == nil {
		opt.Filter = bson.M{}
	}
	if opt.Asc != nil || opt.Desc != nil {
		var sort bson.D
		for _, v := range opt.Asc {
			sort = append(sort, bson.E{Key: v, Value: 1})
		}
		for _, v := range opt.Desc {
			sort = append(sort, bson.E{Key: v, Value: -1})
		}
		f.Sort = sort
	}
	if opt.Select != nil || opt.Exclude != nil {
		var projection bson.D
		for _, v := range opt.Select {
			projection = append(projection, bson.E{Key: v, Value: 1})
		}
		for _, v := range opt.Desc {
			projection = append(projection, bson.E{Key: v, Value: 0})
		}
		f.Projection = projection
	}
	return f
}
