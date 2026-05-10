package mgd

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Id(id any) bson.ObjectID {
	if idStr, ok := id.(string); ok {
		objectId, _ := bson.ObjectIDFromHex(idStr)
		return objectId
	}
	if objectId, ok := id.(bson.ObjectID); ok {
		return objectId
	}
	return bson.NilObjectID
}
func Ids(ids any) interface{} {
	if idStr, ok := ids.([]string); ok {
		var newIds []bson.ObjectID
		for _, v := range idStr {
			objectId, _ := bson.ObjectIDFromHex(v)
			newIds = append(newIds, objectId)
		}
		return newIds
	}
	if objectIds, ok := ids.([]bson.ObjectID); ok {
		return objectIds
	}
	return []bson.ObjectID{}
}

func findOneOptions(opt Example) *options.FindOneOptionsBuilder {
	f := options.FindOne()
	if opt.Asc != nil || opt.Desc != nil {
		var sort bson.D
		for _, v := range opt.Asc {
			sort = append(sort, bson.E{Key: v, Value: 1})
		}
		for _, v := range opt.Desc {
			sort = append(sort, bson.E{Key: v, Value: -1})
		}
		f.SetSort(sort)
	}
	if opt.Select != nil || opt.Exclude != nil || len(opt.Project) > 0 {
		var projection bson.D
		for _, v := range opt.Select {
			projection = append(projection, bson.E{Key: v, Value: 1})
		}
		for _, v := range opt.Exclude {
			projection = append(projection, bson.E{Key: v, Value: 0})
		}
		for k, v := range opt.Project {
			projection = append(projection, bson.E{Key: k, Value: v})
		}
		f.SetProjection(projection)
	}
	return f
}

func findOptions(opt Example) *options.FindOptionsBuilder {
	f := options.Find()
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
		if sort != nil { //如果是空会报错
			f.SetSort(sort)
		}
	}
	if opt.Select != nil || opt.Exclude != nil || len(opt.Project) > 0 {
		var projection bson.D
		for _, v := range opt.Select {
			projection = append(projection, bson.E{Key: v, Value: 1})
		}
		for _, v := range opt.Exclude {
			projection = append(projection, bson.E{Key: v, Value: 0})
		}
		for k, v := range opt.Project {
			projection = append(projection, bson.E{Key: k, Value: v})
		}
		f.SetProjection(projection)
	}
	if opt.BatchSize != 0 {
		f.SetBatchSize(opt.BatchSize)
	}
	return f
}
