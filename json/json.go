package json

import (
	j "encoding/json"
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"unsafe"
)

func Marshal(v interface{}) ([]byte, error) {
	return j.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return j.MarshalIndent(v, prefix, indent)
}

//func Unmarshal(data []byte, v interface{}) error {
//	return j.Unmarshal(data, v)
//}

func NewEncoder(writer io.Writer) *j.Encoder {
	return j.NewEncoder(writer)
}

func NewDecoder(reader io.Reader) *j.Decoder {
	return j.NewDecoder(reader)
}

//-------------custom-----------------

//超长复杂json转map时，json-iterator效率高，简单json转map时gjson效率高
//获取json中的值时gjson速度快，约为json-iterator的10倍

// extra.SetNamingStrategy(LowerCaseWithUnderscores)//统一命名风格
// extra.SupportPrivateFields()//启用私有字段
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	extra.RegisterFuzzyDecoders()
	jsoniter.RegisterTypeDecoderFunc("bool", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		ty := iter.WhatIsNext()
		switch ty {
		case jsoniter.BoolValue:
			*((*bool)(ptr)) = iter.ReadBool()
		case jsoniter.NumberValue:
			number := iter.ReadNumber()
			if n, err := number.Int64(); err == nil {
				if n > 0 {
					*((*bool)(ptr)) = true
				} else {
					*((*bool)(ptr)) = false
				}
			} else {
				*((*bool)(ptr)) = false
			}
		case jsoniter.StringValue:
			str := iter.ReadString()
			if str == "true" {
				*((*bool)(ptr)) = true
			} else {
				*((*bool)(ptr)) = false
			}
		case jsoniter.NilValue:
			iter.ReadNil()
			*((*bool)(ptr)) = false
		default:
			*((*bool)(ptr)) = false
		}
	})
}

type Result struct {
	data gjson.Result
}

func String(value interface{}) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
func Bytes(value interface{}) []byte {
	bytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return bytes
}

//	func Unmarshal(data []byte, v interface{}) error {
//		return j.Unmarshal(data, v)
//	}
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// 临时使用
func ArrayStrToString(value string) string {
	a := []interface{}{}
	_ = json.Unmarshal([]byte(value), &a)
	return String(a)
}

// gjson
func (r Result) String() string {
	return r.data.String()
}
func (r Result) Value() interface{} {
	return r.data.Value()
}
func (r Result) Map() map[string]Result {
	result := map[string]Result{}
	for k, v := range r.data.Map() {
		result[k] = Result{data: v}
	}
	return result
}
func (r Result) MapAny() map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range r.data.Map() {
		result[k] = v.Value()
	}
	return result
}
func (r Result) Get(key string) Result {
	return Result{data: r.data.Get(key)}
}
func (r Result) GetString(key string) string {
	return r.data.Get(key).String()
}
func (r Result) GetInt(key string) int64 {
	return r.data.Get(key).Int()
}
func (r Result) GetInt64(key string) int64 {
	return r.data.Get(key).Int()
}
func (r Result) GetFloat32(key string) float32 {
	return float32(r.data.Get(key).Float())
}
func (r Result) GetFloat64(key string) float64 {
	return r.data.Get(key).Float()
}
func (r Result) GetBool(key string) bool {
	return r.data.Get(key).Bool()
}
func (r Result) GetArray(key string) []Result {
	var array []Result
	for _, v := range r.data.Get(key).Array() {
		array = append(array, Result{data: v})
	}
	return array
}
func (r Result) GetMap(key string) map[string]Result {
	return r.Get(key).Map()
}
func (r Result) GetValue(key string) interface{} {
	return r.data.Get(key).Value()
}
func (r Result) Array() []Result {
	var array []Result
	for _, v := range r.data.Array() {
		array = append(array, Result{data: v})
	}
	return array
}

func Parse(data string) *Result {
	return &Result{data: gjson.Parse(data)}
}

func Valid(data string) bool {
	return gjson.Valid(data)
}

func ValidBytes(data []byte) bool {
	return gjson.ValidBytes(data)
}

func To(str string, v interface{}) {
	err := json.UnmarshalFromString(str, v)
	if err != nil {
		panic(err)
	}
}
func ToWithError(str string, v interface{}) error {
	err := json.UnmarshalFromString(str, v)
	if err != nil {
		return err
	}
	return nil
}
func Set(json, path string, value interface{}) string {
	result, _ := sjson.Set(json, path, value)
	return result
}
func SetRaw(json, path string, value string) string {
	result, _ := sjson.SetRaw(json, path, value)
	return result
}
func Delete(json, path string) string {
	result, _ := sjson.Delete(json, path)
	return result
}
