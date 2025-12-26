package gconv

import (
	"time"

	"github.com/spf13/cast"
)

/************* 字符串相关 *************/
func String(i any) string           { return cast.ToString(i) }
func StringE(i any) (string, error) { return cast.ToStringE(i) }

/************* 布尔值相关 *************/

func Bool(i any) bool                  { return cast.ToBool(i) }
func BoolE(i any) (bool, error)        { return cast.ToBoolE(i) }
func BoolSlice(i any) []bool           { return cast.ToBoolSlice(i) }
func BoolSliceE(i any) ([]bool, error) { return cast.ToBoolSliceE(i) }

/************* 整型相关 *************/
func Int(i any) int                      { return cast.ToInt(i) }
func IntE(i any) (int, error)            { return cast.ToIntE(i) }
func IntSlice(i any) []int               { return cast.ToIntSlice(i) }
func IntSliceE(i any) ([]int, error)     { return cast.ToIntSliceE(i) }
func Int8(i any) int8                    { return cast.ToInt8(i) }
func Int8E(i any) (int8, error)          { return cast.ToInt8E(i) }
func Int8SliceE(i any) ([]int8, error)   { return cast.ToInt8SliceE(i) }
func Int16(i any) int16                  { return cast.ToInt16(i) }
func Int16E(i any) (int16, error)        { return cast.ToInt16E(i) }
func Int16SliceE(i any) ([]int16, error) { return cast.ToInt16SliceE(i) }
func Int32(i any) int32                  { return cast.ToInt32(i) }
func Int32E(i any) (int32, error)        { return cast.ToInt32E(i) }
func Int32SliceE(i any) ([]int32, error) { return cast.ToInt32SliceE(i) }
func Int64(i any) int64                  { return cast.ToInt64(i) }
func Int64E(i any) (int64, error)        { return cast.ToInt64E(i) }
func Int64Slice(i any) []int64           { return cast.ToInt64Slice(i) }
func Int64SliceE(i any) ([]int64, error) { return cast.ToInt64SliceE(i) }

/************* 无符号整型相关 *************/
func Uint(i any) uint                      { return cast.ToUint(i) }
func UintE(i any) (uint, error)            { return cast.ToUintE(i) }
func UintSlice(i any) []uint               { return cast.ToUintSlice(i) }
func UintSliceE(i any) ([]uint, error)     { return cast.ToUintSliceE(i) }
func Uint8(i any) uint8                    { return cast.ToUint8(i) }
func Uint8E(i any) (uint8, error)          { return cast.ToUint8E(i) }
func Uint8SliceE(i any) ([]uint8, error)   { return cast.ToUint8SliceE(i) }
func Uint16(i any) uint16                  { return cast.ToUint16(i) }
func Uint16E(i any) (uint16, error)        { return cast.ToUint16E(i) }
func Uint16SliceE(i any) ([]uint16, error) { return cast.ToUint16SliceE(i) }
func Uint32(i any) uint32                  { return cast.ToUint32(i) }
func Uint32E(i any) (uint32, error)        { return cast.ToUint32E(i) }
func Uint32SliceE(i any) ([]uint32, error) { return cast.ToUint32SliceE(i) }
func Uint64(i any) uint64                  { return cast.ToUint64(i) }
func Uint64E(i any) (uint64, error)        { return cast.ToUint64E(i) }
func Uint64SliceE(i any) ([]uint64, error) { return cast.ToUint64SliceE(i) }

/************* 浮点相关 *************/
func Float32(i any) float32                  { return cast.ToFloat32(i) }
func Float32E(i any) (float32, error)        { return cast.ToFloat32E(i) }
func Float32SliceE(i any) ([]float32, error) { return cast.ToFloat32SliceE(i) }
func Float64(i any) float64                  { return cast.ToFloat64(i) }
func Float64E(i any) (float64, error)        { return cast.ToFloat64E(i) }
func Float64Slice(i any) []float64           { return cast.ToFloat64Slice(i) }
func Float64SliceE(i any) ([]float64, error) { return cast.ToFloat64SliceE(i) }

/************* 时间相关 *************/
func Duration(i any) time.Duration                  { return cast.ToDuration(i) }
func DurationE(i any) (time.Duration, error)        { return cast.ToDurationE(i) }
func DurationSlice(i any) []time.Duration           { return cast.ToDurationSlice(i) }
func DurationSliceE(i any) ([]time.Duration, error) { return cast.ToDurationSliceE(i) }
func Time(i any) time.Time                          { return cast.ToTime(i) }
func TimeE(i any) (time.Time, error)                { return cast.ToTimeE(i) }
func TimeInDefaultLocation(i any, location *time.Location) time.Time {
	return cast.ToTimeInDefaultLocation(i, location)
}
func TimeInDefaultLocationE(i any, location *time.Location) (time.Time, error) {
	return cast.ToTimeInDefaultLocationE(i, location)
}

/************* 通用 slice / map 转换 *************/
func Slice(i any) []any                        { return cast.ToSlice(i) }
func SliceE(i any) ([]any, error)              { return cast.ToSliceE(i) }
func StringSlice(i any) []string               { return cast.ToStringSlice(i) }
func StringSliceE(i any) ([]string, error)     { return cast.ToStringSliceE(i) }
func StringMap(i any) map[string]any           { return cast.ToStringMap(i) }
func StringMapE(i any) (map[string]any, error) { return cast.ToStringMapE(i) }
func StringMapBool(i any) map[string]bool      { return cast.ToStringMapBool(i) }
func StringMapBoolE(i any) (map[string]bool, error) {
	return cast.ToStringMapBoolE(i)
}
func StringMapInt(i any) map[string]int { return cast.ToStringMapInt(i) }
func StringMapIntE(i any) (map[string]int, error) {
	return cast.ToStringMapIntE(i)
}
func StringMapInt64(i any) map[string]int64 { return cast.ToStringMapInt64(i) }
func StringMapInt64E(i any) (map[string]int64, error) {
	return cast.ToStringMapInt64E(i)
}
func StringMapString(i any) map[string]string { return cast.ToStringMapString(i) }
func StringMapStringE(i any) (map[string]string, error) {
	return cast.ToStringMapStringE(i)
}
func StringMapStringSlice(i any) map[string][]string {
	return cast.ToStringMapStringSlice(i)
}
func StringMapStringSliceE(i any) (map[string][]string, error) {
	return cast.ToStringMapStringSliceE(i)
}

//======time相关======
//func StringToDateInDefaultLocation(s string, location *time.Location) time.Time {
//	return cast.StringToDateInDefaultLocation(s, location)
//}
//2）map ↔ struct 转换
//
//✅ 第一推荐：github.com/go-viper/mapstructure/v2
//（而不是已归档的 mitchellh/mapstructure）

//基础类型转换：
//
//引 github.com/spf13/cast
//
//map ↔ struct：
//
//用 github.com/go-viper/mapstructure/v2
//
//struct ↔ struct：
//
//用 github.com/jinzhu/copier
//
//配置合并：
//
//用 dario.cat/mergo（如果有配置层需求）
