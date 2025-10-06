package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/basicfu/gf/g"
	"github.com/basicfu/gf/util/gconv"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

// 从redigo更新到go-redis封装更好
func Init(redisURI string) {
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		panic(err)
	}
	rdb = redis.NewClient(opt)
}

var SetLog = redis.SetLogger

//	func SetLog() {
//		redis.SetLogger()
//	}
type Item struct {
	Id   string
	Data map[string]string //stream key不重复不为空
}
type Result struct {
	val any
}

func New[T any](cmd interface {
	Err() error
	Val() T
}) Result {
	_panic(cmd.Err())
	return Result{val: cmd.Val()}
}
func (r Result) Data() any {
	return r.val
}
func (r Result) Bool() bool {
	return gconv.Bool(r.val)
}
func (r Result) String() string { //不正确应该
	return gconv.String(r.val)
}
func (r Result) Strings() []string {
	return gconv.SliceStr(r.val)
}
func (r Result) Int() int {
	return gconv.Int(r.val)
}
func (r Result) Int64() int64 {
	return gconv.Int64(r.val)
}
func (r Result) MapStringString() map[string]string {
	return gconv.MapStrStr(r.val)
}

//
//// 获取第一个key，第一个id，并转为map，key重复会覆盖,并过滤空key
//func (r Result) StreamMap() (string, map[string]string) {
//	id := ""
//	if r.cmd == nil {
//		return id, map[string]string{}
//	}
//	vs, _ := red.Values(r.cmd, r.Error)
//	vs, _ = red.Values(vs[0], nil)
//	vs, _ = red.Values(vs[1], nil)
//	vs, _ = red.Values(vs[0], nil)
//	id = string(vs[0].([]byte))
//	vs, _ = red.Values(vs[1], nil)
//	data := map[string]string{}
//	for i := 0; i < len(vs)/2; i++ {
//		key := string(vs[i*2].([]byte))
//		if key != "" {
//			data[key] = string(vs[i*2+1].([]byte))
//		}
//	}
//	return id, data
//}
//
//// maps用于一次拉取多个元素的结果，多个id对应map
//func (r Result) StreamMaps() []Item {
//	items := []Item{}
//	if r.cmd == nil {
//		return items
//	}
//	vs, _ := red.Values(r.cmd, r.Error)
//	vs, _ = red.Values(vs[0], nil)
//	vs, _ = red.Values(vs[1], nil)
//	for _, v := range vs {
//		item, _ := red.Values(v, nil)
//		id := string(item[0].([]byte))
//		values, _ := red.Values(item[1], nil)
//		data := map[string]string{}
//		for i := 0; i < len(values)/2; i++ {
//			key := string(values[i*2].([]byte))
//			if key != "" {
//				data[key] = string(values[i*2+1].([]byte))
//			}
//		}
//		items = append(items, Item{
//			Id:   id,
//			Data: data,
//		})
//	}
//	return items
//}

// 封装的组件要尽量独立减少使用其他封装的依赖
func _panic(error error) {
	if error != nil {
		if !errors.Is(error, redis.Nil) {
			panic(error.Error())
		}
	}
}

// ===========pub/sub========
func Publish(channel string, message any) Result {
	return New[int64](rdb.Publish(ctx, channel, message))
}

type PubSub = redis.PubSub

func Subscribe(channel string) PubSub {
	return *rdb.Subscribe(ctx, channel)
}

// ===========lua========
type Script struct {
	s *redis.Script
}

func NewScript(src string) Script {
	return Script{
		s: redis.NewScript(src),
	}
}
func (s Script) Run(keys []string, args ...interface{}) Result {
	return New[any](s.s.Run(ctx, rdb, keys, args...))
}

// =========cmd===========
func Set(key string, value interface{}) Result {
	return New[string](rdb.Set(ctx, key, value, 0))
}
func SetEx(key string, value interface{}, expiration time.Duration) Result {
	return New[string](rdb.SetEx(ctx, key, value, expiration))
}
func SetNx(key string, value interface{}) Result {
	return New[bool](rdb.SetNX(ctx, key, value, 0))
}
func SetNxEx(key string, value interface{}, expiration time.Duration) Result {
	return New[bool](rdb.SetNX(ctx, key, value, expiration))
}

func Get(key string) Result {
	return New[string](rdb.Get(ctx, key))
}
func Del(keys ...string) int64 {
	return New[int64](rdb.Del(ctx, keys...)).Int64()
}
func Keys(pattern string) []string {
	return New[[]string](rdb.Keys(ctx, pattern)).Strings()
}
func Exists(key string) bool {
	val, err := rdb.Exists(ctx, key).Result()
	_panic(err)
	return val != 0
}
func ExistsCount(key ...string) int64 {
	val, err := rdb.Exists(ctx, key...).Result()
	_panic(err)
	return val
}

// 返回每个key是否存在
func ExistsBatch(keys ...string) g.MapStrBool {
	cmds, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.Exists(ctx, key)
		}
		return nil
	})
	_panic(err)
	for _, cmd := range cmds {
		fmt.Println(cmd.(*redis.StringCmd).Val())
	}
	result := g.MapStrBool{}
	for index, cmd := range cmds {
		result[keys[index]] = cmd.(*redis.IntCmd).Val() == 1
	}
	return result
}

// ===========hash========
func HSet(key string, values ...any) Result {
	return New[int64](rdb.HSet(ctx, key, values...))
}

// 设置hash key的过期时间ms
func HSetEx(key string, exMs int64, values ...any) int64 {
	return hSetExWithArgs(key, redis.HSetEXOptions{
		Condition:      "",
		ExpirationType: redis.HSetEXExpirationPX, //ms，还可加PXAT指定unix时间
		ExpirationVal:  exMs,
	}, values...)
}

// 同时附加FNX，返回1设置成功，0数据存在不设置
func HSetExNx(key string, exMs int64, values ...any) int64 {
	return hSetExWithArgs(key, redis.HSetEXOptions{
		Condition:      redis.HSetEXFNX,
		ExpirationType: redis.HSetEXExpirationPX, //ms
		ExpirationVal:  exMs,
	}, values...)
}

// 同时附加FXX
func HSetExXx(key string, exMs int64, values ...any) int64 {
	return hSetExWithArgs(key, redis.HSetEXOptions{
		Condition:      redis.HSetEXFXX,
		ExpirationType: redis.HSetEXExpirationPX, //ms
		ExpirationVal:  exMs,
	}, values...)
}
func hSetExWithArgs(key string, opt redis.HSetEXOptions, values ...any) int64 {
	//底层库拼接是支持的，但是传参确是string，这里转换下使上层支持any
	var str []string
	for _, v := range values {
		str = append(str, fmt.Sprint(v))
	}
	val, err := rdb.HSetEXWithArgs(ctx, key, &opt, str...).Result()
	_panic(err)
	return val
}

func HGet(key, field string) Result {
	return New[string](rdb.HGet(ctx, key, field))
}

func HMGet(key string, fields ...string) Result {
	return New[[]any](rdb.HMGet(ctx, key, fields...))
}
func HDel(key string, fields ...string) Result {
	return New[int64](rdb.HDel(ctx, key, fields...))
}
func HGetAll(key string) Result {
	return New[map[string]string](rdb.HGetAll(ctx, key))
}

// ===========stream========
type XMessage = redis.XMessage

func XAdd(stream string, value any) Result {
	return New[string](rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: value,
	}))
}
func XAddWithId(stream, id string, value any) Result {
	return New[string](rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		ID:     id,
		Values: value,
	}))
}

// 返回true已存在
func XGroupCreateMkStream(stream, group, start string) bool {
	_, err := rdb.XGroupCreateMkStream(ctx, stream, group, start).Result()
	if err != nil {
		if strings.Contains(err.Error(), "BUSYGROUP") {
			return true
		}
		_panic(err)
	}
	return false
}

// 仅读取单个stream，返回结果体也简单一些
func XReadGroup(group, consumer string, count int64, block time.Duration, stream, id string) []XMessage {
	val, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, id}, //[e.g. stream1 stream2 id1 id2]
		Count:    count,
		Block:    block,
	}).Result()
	_panic(err)
	if len(val) == 1 {
		return val[0].Messages
	}
	return []XMessage{}
}

type XAutoClaimArgs = redis.XAutoClaimArgs

func XAutoClaim(args XAutoClaimArgs) ([]XMessage, string) {
	val, start, err := rdb.XAutoClaim(ctx, &args).Result()
	_panic(err)
	return val, start
}

// 使用默认mode
func XAckDel(stream, group string, ids ...string) []interface{} {
	val, err := rdb.XAckDel(ctx, stream, group, "KEEPREF", ids...).Result()
	_panic(err)
	return val
}

//
//func SPop(key interface{}, count int64) []string {
//	s, _ := red.Strings(exec("spop", key, count), nil)
//	return s
//}
//func LPop(key interface{}) string {
//	s, _ := red.String(exec("lpop", key), nil)
//	return s
//}
//func LLen(key interface{}) int64 {
//	s, _ := red.Int64(exec("llen", key), nil)
//	return s
//}
//func SRandMember(key interface{}, count int64) []string {
//	s, _ := red.Strings(exec("SRANDMEMBER", key, count), nil)
//	return s
//}
//func SRem(key interface{}, values ...string) int64 {
//	args := []interface{}{key}
//	for _, val := range values {
//		args = append(args, val)
//	}
//	i, _ := red.Int64(exec("srem", args...), nil)
//	return i
//}
//func SCard(key interface{}) int64 {
//	n, _ := red.Int64(exec("SCARD", key), nil)
//	return n
//}
//
//// 已默认添加目标key
//func SUnionStore(sourceKey interface{}, targetKey interface{}) int64 {
//	s, _ := red.Int64(exec("SUNIONSTORE", targetKey, targetKey, sourceKey, "temp"), nil)
//	return s
//}
//func Lpush(key interface{}, values ...interface{}) {
//	args := []interface{}{key}
//	for _, val := range values {
//		args = append(args, val)
//	}
//	i := exec("lpush", args...)
//	_, _ = red.Strings(i, nil)
//}
//func RPop(key interface{}) string {
//	s, _ := red.String(exec("rpop", key), nil)
//	return s
//}
//func Rpush(key interface{}, values ...interface{}) {
//	args := []interface{}{key}
//	for _, val := range values {
//		args = append(args, val)
//	}
//	i := exec("rpush", args...)
//	_, _ = red.Strings(i, nil)
//}
//func SAdd(key interface{}, values ...string) int64 {
//	args := []interface{}{key}
//	for _, val := range values {
//		args = append(args, val)
//	}
//	i, _ := red.Int64(exec("sadd", args...), nil)
//	return i
//}
//func Ttl(key interface{}) int64 {
//	s, _ := red.Int64(exec("ttl", key), nil)
//	return s
//}
//func Expire(key interface{}, time int64) {
//	_, _ = red.Int64(exec("expire", key, time), nil)
//}
//func SetEx(key interface{}, value interface{}, time int64) {
//	_, _ = red.Int64(exec("setex", key, time, value), nil)
//}
//func SetExNx(key interface{}, value interface{}, time int64) bool {
//	res := exec("set", key, value, "ex", time, "nx")
//	if res == nil {
//		return false
//	}
//	//res string Ok
//	return true
//}
//func LRange(key interface{}, startIndex int, endIndex int) []string {
//	strings, _ := red.Strings(exec("lrange", key, startIndex, endIndex), nil)
//	return strings
//}
//func LRem(key interface{}, count int, value interface{}) Result {
//	return Result{cmd: exec("lrem", key, count, value)}
//}
//func Ltrim(key interface{}, startIndex int, endIndex int) {
//	_, _ = red.Strings(exec("ltrim", key, startIndex, endIndex), nil)
//}
//func HExists(key interface{}, hk interface{}) bool {
//	s, _ := red.Bool(exec("hexists", key, hk), nil)
//	return s
//}
//func SisMember(key interface{}, value interface{}) bool {
//	s, _ := red.Bool(exec("SISMEMBER", key, value), nil)
//	return s
//}
//func SMembers(key interface{}) []string {
//	i := exec("SMEMBERS", key)
//	s, _ := red.Strings(i, nil)
//	return s
//}
//func HIncrBy(key interface{}, hk interface{}, incr int64) Result {
//	return Result{cmd: exec("hincrby", key, hk, incr)}
//}
//func IncrBy(key interface{}, incr int64) Result {
//	return Result{cmd: exec("incrby", key, incr)}
//}
//func Incr(key interface{}) Result {
//	return Result{cmd: exec("incr", key)}
//}

//func ZAdd(key interface{}, score interface{}, value interface{}) Result {
//	return Result{cmd: exec("zadd", key, score, value)}
//}
//func ZRangeByScore(key interface{}, min, max interface{}) Result {
//	return Result{cmd: exec("ZRANGEBYSCORE", key, min, max, "WITHSCORES")} //key score
//}
//func ZRevRangeByScore(key interface{}, max, min int, params ...interface{}) Result {
//	return Result{cmd: exec("ZREVRANGEBYSCORE", key, max, min, "WITHSCORES")} //key score
//}
//func ZIncrBy(key interface{}, incr int, member interface{}) Result {
//	return Result{cmd: exec("ZINCRBY", key, incr, member)}
//}
//func ZRange(key interface{}, min, max int, params ...interface{}) Result {
//	return Result{cmd: exec("ZRANGE", key, min, max)}
//}
//func ZRem(key interface{}, value interface{}) Result {
//	return Result{cmd: exec("ZREM", key, value)}
//}
//func ZAddMany(key interface{}, array ...[]interface{}) Result {
//	args := []interface{}{key}
//	for _, val := range array {
//		args = append(args, val[0], val[1])
//	}
//	return Result{cmd: exec("zadd", args...)}
//}
//func Exec(cmd string, values ...interface{}) Result {
//	return Result{cmd: exec(cmd, values...)}
//}
//

//func execWithTimeout(cmd string, args ...interface{}) (interface{}, error) {
//	if pool == nil {
//		return nil, errors.New("")
//	}
//	con := pool.Get()
//	if err := con.Err(); err != nil {
//		return nil, err
//	}
//	defer con.Close()
//	return red.DoWithTimeout(con, time.Duration(60*1000)*time.Millisecond, cmd, args...)
//}
