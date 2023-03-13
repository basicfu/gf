package mgd

import (
	"errors"
	"github.com/basicfu/gf/g"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

func Id(id any) primitive.ObjectID {
	if idStr, ok := id.(string); ok {
		objectId, _ := primitive.ObjectIDFromHex(idStr)
		return objectId
	}
	if objectId, ok := id.(primitive.ObjectID); ok {
		return objectId
	}
	return primitive.NilObjectID
}
func Ids(ids any) interface{} {
	if idStr, ok := ids.([]string); ok {
		var newIds []primitive.ObjectID
		for _, v := range idStr {
			objectId, _ := primitive.ObjectIDFromHex(v)
			newIds = append(newIds, objectId)
		}
		return newIds
	}
	if objectIds, ok := ids.([]primitive.ObjectID); ok {
		return objectIds
	}
	return []primitive.ObjectID{}
}

//struct过滤0值，map不过滤，其他类型考虑是否抛错
func toFilter(params any) g.Map {
	filter := g.Map{}
	if params == nil {
		return filter
	}
	t := reflect.TypeOf(params)
	kind := t.Kind()
	if kind == reflect.Struct {
		v := reflect.ValueOf(params)
		for i := 0; i < t.NumField(); i++ {
			value := v.Field(i)
			if !value.IsZero() {
				filter[t.Field(i).Tag.Get("bson")] = value.Interface()
			}
		}
	} else if kind == reflect.Map {
		switch params.(type) {
		case map[string]interface{}:
			filter = params.(map[string]interface{})
		default:
			panic(errors.New("不支持的类型"))
		}
	} else {
		panic(errors.New("不支持的类型"))
	}
	return filter
}
func findOneOptions(opt Example) options.FindOneOptions {
	f := options.FindOneOptions{}
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
		f.Projection = projection
	}
	return f
}

func findOptions(opt Example) options.FindOptions {
	f := options.FindOptions{}
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
		for _, v := range opt.Exclude {
			projection = append(projection, bson.E{Key: v, Value: 0})
		}
		f.Projection = projection
	}
	if opt.BatchSize != 0 {
		f.SetBatchSize(opt.BatchSize)
	}
	return f
}
