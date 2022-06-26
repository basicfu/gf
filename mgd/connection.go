package mgd

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var config *Config
var client *mongo.Client
var db *mongo.Database

type Config struct {
	CtxTimeout time.Duration
}

func ctx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), config.CtxTimeout)
	//ctx = context.Background()
	return ctx
}
func Init(conf *Config, dbName string, opts ...*options.ClientOptions) {
	if conf == nil {
		conf = &Config{CtxTimeout: 10 * time.Second}
	}
	config = conf
	var err error
	client, err = mongo.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	if err = client.Connect(ctx()); err != nil {
		panic(err)
	}
	db = client.Database(dbName)
}

func Coll(m interface{}, opts ...*options.CollectionOptions) *Collection {
	name := ""
	if collNameGetter, ok := m.(CollectionNameGetter); ok {
		name = collNameGetter.CollectionName()
	} else {
		name = reflect.TypeOf(m).Elem().Name()
	}
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(name, "${1}_${2}")
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	name = strings.ToLower(snake)
	coll := db.Collection(name, opts...)
	return &Collection{coll: coll, model: m}
}
