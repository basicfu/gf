package redis

import (
	"errors"
	"github.com/basicfu/gf/g"
	"github.com/basicfu/gf/util/gconv"
	red "github.com/gomodule/redigo/redis"
	"log"
	"time"
)

var pool *red.Pool

//选择该库因为比较简介，自己拼接命令执行do即可
func Init(address string, password string) {
	pool = &red.Pool{
		MaxIdle:     256,
		MaxActive:   0,
		IdleTimeout: 10 * time.Second,
		Dial: func() (red.Conn, error) {
			return red.Dial(
				"tcp",
				address,
				red.DialReadTimeout(60*time.Second+10*time.Second),
				red.DialWriteTimeout(10*time.Second),
				red.DialConnectTimeout(10*time.Second),
				red.DialDatabase(0),
				red.DialPassword(password),
			)
		},
	}
}
func InitWithDb(address string, password string, db int) {
	pool = &red.Pool{
		MaxIdle:     256,
		MaxActive:   0,
		IdleTimeout: 10 * time.Second,
		Dial: func() (red.Conn, error) {
			return red.Dial(
				"tcp",
				address,
				red.DialReadTimeout(60*time.Second+10*time.Second),
				red.DialWriteTimeout(10*time.Second),
				red.DialConnectTimeout(10*time.Second),
				red.DialDatabase(db),
				red.DialPassword(password),
			)
		},
	}
}

type Result struct {
	data  interface{}
	Error error
}
type Item struct {
	Id   string
	Data map[string]string //stream key不重复不为空
}

func (r Result) IsNil() bool {
	return r.data == nil
}
func (r Result) Bool() bool {
	if r.data == nil {
		return false
	}
	s, err := red.Bool(r.data, nil)
	_panic(err)
	return s
}
func (r Result) String() string {
	if r.data == nil {
		return ""
	}
	s, _ := red.String(r.data, nil)
	return s
}
func (r Result) Strings() []string {
	if r.data == nil {
		return []string{}
	}
	s, err := red.Strings(r.data, nil)
	_panic(err)
	return s
}
func (r Result) Int() int {
	s, _ := red.Int(r.data, nil)
	return s
}
func (r Result) Int64() int64 {
	s, _ := red.Int64(r.data, nil)
	return s
}
func (r Result) Map() map[string]string {
	s, err := red.StringMap(r.data, nil)
	_panic(err)
	return s
}

//获取第一个key，第一个id，并转为map，key重复会覆盖,并过滤空key
func (r Result) StreamMap() (string, map[string]string) {
	id := ""
	if r.data == nil {
		return id, map[string]string{}
	}
	vs, _ := red.Values(r.data, r.Error)
	vs, _ = red.Values(vs[0], nil)
	vs, _ = red.Values(vs[1], nil)
	vs, _ = red.Values(vs[0], nil)
	id = string(vs[0].([]byte))
	vs, _ = red.Values(vs[1], nil)
	data := map[string]string{}
	for i := 0; i < len(vs)/2; i++ {
		key := string(vs[i*2].([]byte))
		if key != "" {
			data[key] = string(vs[i*2+1].([]byte))
		}
	}
	return id, data
}

//maps用于一次拉取多个元素的结果，多个id对应map
func (r Result) StreamMaps() []Item {
	items := []Item{}
	if r.data == nil {
		return items
	}
	vs, _ := red.Values(r.data, r.Error)
	vs, _ = red.Values(vs[0], nil)
	vs, _ = red.Values(vs[1], nil)
	for _, v := range vs {
		item, _ := red.Values(v, nil)
		id := string(item[0].([]byte))
		values, _ := red.Values(item[1], nil)
		data := map[string]string{}
		for i := 0; i < len(values)/2; i++ {
			key := string(values[i*2].([]byte))
			if key != "" {
				data[key] = string(values[i*2+1].([]byte))
			}
		}
		items = append(items, Item{
			Id:   id,
			Data: data,
		})
	}
	return items
}

//封装的组件要尽量独立减少使用其他封装的依赖
func _panic(error error) {
	if error != nil {
		panic(error.Error())
	}
}
func Publish(channel string, message string) {
	exec("PUBLISH", channel, message)
}

func Subscribe(channel string, messageFunc func(data string), subscriptionFunc func(kind string)) {
	run := func() {
		con := pool.Get()
		defer func() {
			_ = con.Close()
		}()
		if con.Err() != nil {
			panic(con.Err().Error())
		}
		psc := red.PubSubConn{Conn: con}
		err := psc.Subscribe(channel)
		defer func() {
			_ = psc.Unsubscribe()
		}()
		if err != nil {
			panic(err.Error())
		}
		done := make(chan error, 1)
		ticker := time.NewTicker(10 * time.Second)
		defer close(done) //关闭，否则可能会触发多次失败
		defer ticker.Stop()
		go func() {
			for {
				switch v := psc.ReceiveWithTimeout(60 * time.Second).(type) {
				case error:
					log.Println("redis获取消息失败", v.Error())
					done <- v
					return
				case red.Message:
					messageFunc(string(v.Data))
				case red.Subscription:
					if subscriptionFunc != nil {
						subscriptionFunc(v.Kind)
					}
				}
			}
		}()
		for {
			select {
			case <-ticker.C:
				if err = psc.Ping(""); err != nil {
					return
				}
			case _ = <-done:
				return
			}
		}
	}
	for {
		g.TryBlock(func() {
			run()
		}, func(err error) {
			log.Println("redis订阅失败", err.Error())
		})
		time.Sleep(1 * time.Second)
	}
}
func exec(cmd string, args ...interface{}) interface{} {
	con := pool.Get()
	_panic(con.Err())
	defer con.Close()
	do, err := con.Do(cmd, red.Args{}.AddFlat(args)...)
	_panic(err)
	return do
}
func NewScript(keyCount int, src string) *red.Script {
	return red.NewScript(keyCount, src)
}
func Lua(script *red.Script, args ...interface{}) interface{} {
	con := pool.Get()
	if err := con.Err(); err != nil {
		panic(err.Error())
	}
	defer con.Close()
	do, err := script.Do(con, args...)
	_panic(err)
	return do
}

func LrangeAndLtrim(key string, startLrangeIndex int, endLrangeIndex int, startTrimIndex int, endTrimIndex int) []string {
	con := pool.Get()
	if err := con.Err(); err != nil {
		return nil
	}
	defer con.Close()
	_ = con.Send("multi")
	_ = con.Send("lrange", key, startLrangeIndex, endLrangeIndex)
	_ = con.Send("ltrim", key, startTrimIndex, endTrimIndex)
	r, err := red.Values(con.Do("exec"))
	if err != nil {
		panic("LrangeAndLtrim发生错误：" + err.Error())
	}
	if r[0] == nil {
		return []string{}
	}
	var items []string
	for _, v := range r[0].([]interface{}) {
		s := string(v.([]uint8))
		items = append(items, s)
	}
	return items
}
func Del(key interface{}) {
	_ = exec("del", key)
}
func Keys(key interface{}) []string {
	s, _ := red.Strings(exec("keys", key), nil)
	return s
}
func Set(key interface{}, value interface{}) {
	_ = exec("set", key, value)
}
func GetString(key interface{}) string {
	s, _ := red.String(exec("get", key), nil)
	return s
}
func Get(key interface{}) Result {
	return Result{data: exec("get", key)}
}

//func BRPopString(key interface{},timeout time.Duration) string {
//	s, err := red.String(execWithTimeout(timeout, "brpop", key, timeout))
//	println(err.Error())
//	return s
//	s, _ := red.Strings(exec("spop", key, count), nil)
//	return s
//}

func SPop(key interface{}, count int64) []string {
	s, _ := red.Strings(exec("spop", key, count), nil)
	return s
}
func LPop(key interface{}) string {
	s, _ := red.String(exec("lpop", key), nil)
	return s
}
func LLen(key interface{}) int64 {
	s, _ := red.Int64(exec("llen", key), nil)
	return s
}
func SRandMember(key interface{}, count int64) []string {
	s, _ := red.Strings(exec("SRANDMEMBER", key, count), nil)
	return s
}
func SRem(key interface{}, value interface{}) int64 {
	n, _ := red.Int64(exec("SREM", key, value), nil)
	return n
}
func SCard(key interface{}) int64 {
	n, _ := red.Int64(exec("SCARD", key), nil)
	return n
}

//已默认添加目标key
func SUnionStore(sourceKey interface{}, targetKey interface{}) int64 {
	s, _ := red.Int64(exec("SUNIONSTORE", targetKey, targetKey, sourceKey, "temp"), nil)
	return s
}
func Lpush(key interface{}, values ...interface{}) {
	args := []interface{}{key}
	for _, val := range values {
		args = append(args, val)
	}
	i := exec("lpush", args...)
	_, _ = red.Strings(i, nil)
}
func RPop(key interface{}) string {
	s, _ := red.String(exec("rpop", key), nil)
	return s
}
func Rpush(key interface{}, values ...interface{}) {
	args := []interface{}{key}
	for _, val := range values {
		args = append(args, val)
	}
	i := exec("rpush", args...)
	_, _ = red.Strings(i, nil)
}
func SAdd(key interface{}, values ...string) int64 {
	args := []interface{}{key}
	for _, val := range values {
		args = append(args, val)
	}
	i, _ := red.Int64(exec("sadd", args...), nil)
	return i
}
func Ttl(key interface{}) int64 {
	s, _ := red.Int64(exec("ttl", key), nil)
	return s
}
func Expire(key interface{}, time int64) {
	_, _ = red.Int64(exec("expire", key, time), nil)
}
func SetEx(key interface{}, value interface{}, time int64) {
	_, _ = red.Int64(exec("setex", key, time, value), nil)
}
func LRange(key interface{}, startIndex int, endIndex int) []string {
	strings, _ := red.Strings(exec("lrange", key, startIndex, endIndex), nil)
	return strings
}
func LRem(key interface{}, count int, value interface{}) Result {
	return Result{data: exec("lrem", key, count, value)}
}
func Ltrim(key interface{}, startIndex int, endIndex int) {
	_, _ = red.Strings(exec("ltrim", key, startIndex, endIndex), nil)
}
func HExists(key interface{}, hk interface{}) bool {
	s, _ := red.Bool(exec("hexists", key, hk), nil)
	return s
}
func SisMember(key interface{}, value interface{}) bool {
	s, _ := red.Bool(exec("SISMEMBER", key, value), nil)
	return s
}
func SMembers(key interface{}) []string {
	i := exec("SMEMBERS", key)
	s, _ := red.Strings(i, nil)
	return s
}
func HIncrBy(key interface{}, hk interface{}, incr int64) Result {
	return Result{data: exec("hincrby", key, hk, incr)}
}
func IncrBy(key interface{}, incr int64) Result {
	return Result{data: exec("incrby", key, incr)}
}
func Incr(key interface{}) Result {
	return Result{data: exec("incr", key)}
}
func HSet(key interface{}, hk interface{}, hv interface{}) bool {
	s, _ := red.Bool(exec("hset", key, hk, hv), nil)
	return s
}

func HSetMap(key interface{}, hash map[interface{}]interface{}) {
	args := []interface{}{key}
	for k, v := range hash {
		args = append(args, k)
		args = append(args, v)
	}
	_, _ = red.Bool(exec("hset", args...), nil)
}

func HGet(key interface{}, hk interface{}) Result {
	return Result{data: exec("hget", key, hk)}
}
func HMGet(key interface{}, hk ...string) Result {
	return Result{data: exec("hmget", append([]interface{}{key}, gconv.Interfaces(hk)...)...)}
}
func HDel(key interface{}, hk interface{}) Result {
	return Result{data: exec("hdel", key, hk)}
}
func HGetAll(key interface{}) Result {
	return Result{data: exec("hgetall", key)}
}
func ZAdd(key interface{}, score interface{}, value interface{}) Result {
	return Result{data: exec("zadd", key, score, value)}
}
func ZRangeByScore(key interface{}, min, max interface{}) Result {
	return Result{data: exec("ZRANGEBYSCORE", key, min, max, "WITHSCORES")} //key score
}
func ZRevRangeByScore(key interface{}, max, min int, params ...interface{}) Result {
	return Result{data: exec("ZREVRANGEBYSCORE", key, max, min, "WITHSCORES")} //key score
}
func ZIncrBy(key interface{}, incr int, member interface{}) Result {
	return Result{data: exec("ZINCRBY", key, incr, member)}
}
func ZRange(key interface{}, min, max int, params ...interface{}) Result {
	return Result{data: exec("ZRANGE", key, min, max)}
}
func ZRem(key interface{}, value interface{}) Result {
	return Result{data: exec("ZREM", key, value)}
}
func ZAddMany(key interface{}, array ...[]interface{}) Result {
	args := []interface{}{key}
	for _, val := range array {
		args = append(args, val[0], val[1])
	}
	return Result{data: exec("zadd", args...)}
}
func Exec(cmd string, values ...interface{}) Result {
	return Result{data: exec(cmd, values...)}
}
func Exists(key string) bool {
	b, _ := red.Bool(exec("EXISTS", key), nil)
	return b
}

//--stream--
//xadd not-exists-stream nomkstream * username lisi 不自动创建流
//id需要为[整数-整数]格式 field value [field value ...]
func XAdd(key, id string, obj map[string]interface{}) Result {
	args := []interface{}{key, id}
	for k, v := range obj {
		args = append(args, k, v)
	}
	data, err := execWithTimeout("XADD", args...)
	return Result{data: data, Error: err}
}

//block//ms
func XRead(count, block int, key string, id string) Result {
	data, err := execWithTimeout("XREAD", "count", count, "block", block, "streams", key, id)
	return Result{data: data, Error: err}
}

func XDel(key string, id ...string) Result {
	args := []interface{}{key}
	for _, v := range id {
		args = append(args, v)
	}
	return Result{data: exec("XDEL", args...)}
}
func execWithTimeout(cmd string, args ...interface{}) (interface{}, error) {
	if pool == nil {
		return nil, errors.New("")
	}
	con := pool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer con.Close()
	return red.DoWithTimeout(con, time.Duration(60*1000)*time.Millisecond, cmd, args...)
}
