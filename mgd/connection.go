package mgd

import (
	"context"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/basicfu/gf/g"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var config *Config
var client *mongo.Client
var db *mongo.Database

type Config struct {
	CtxTimeout time.Duration
}

func buildCtx(ctxArray ...context.Context) context.Context {
	if len(ctxArray) != 0 {
		return ctxArray[0]
	}
	ctx, _ := context.WithTimeout(context.Background(), config.CtxTimeout)
	return ctx
}
func Init(conf *Config, dbName string, opts ...*options.ClientOptions) {
	if conf == nil {
		conf = &Config{CtxTimeout: 10 * time.Second}
	}
	config = conf
	reg := bson.NewRegistry()
	reg.RegisterTypeDecoder(reflect.TypeOf(decimal.Decimal{}), Decimal{})
	reg.RegisterTypeEncoder(reflect.TypeOf(decimal.Decimal{}), Decimal{})
	decimalOpt := options.Client().SetRegistry(reg)
	allOpts := append(opts, decimalOpt)
	var err error
	client, err = mongo.Connect(allOpts...)
	if err != nil {
		panic(err)
	}
	if err := client.Ping(buildCtx(), nil); err != nil {
		panic(err)
	}
	db = client.Database(dbName)
}
func Close() {
	_ = client.Disconnect(buildCtx())
}
func Coll[T any | g.Map](m *T, opts ...options.Lister[options.CollectionOptions]) *Collection[T] {
	name := reflect.TypeOf(m).Elem().Name()
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(name, "${1}_${2}")
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	name = strings.ToLower(snake)
	coll := db.Collection(name, opts...)
	return &Collection[T]{coll: coll, model: *m}
}
