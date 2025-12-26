package gtime

import (
	"errors"
	"os"
	"strconv"
	"time"
)

func SetTimeZone(zone ...string) error {
	var location *time.Location
	if len(zone) == 0 {
		location = time.FixedZone("CST", 8*3600) // 东八
	} else {
		var err error
		location, err = time.LoadLocation(zone[0])
		if err != nil {
			return err
		}
	}
	time.Local = location
	_ = os.Setenv("TZ", location.String())
	return nil
}

func Unix() int64 {
	return Now().Unix()
}
func UnixMilli() int64 {
	return Now().UnixMilli()
}
func UnixMicro() int64 {
	return Now().UnixMicro()
}
func UnixNano() int64 {
	return Now().UnixNano()
}
func UnixStr() string {
	return Now().UnixStr()
}
func UnixMilliStr() string {
	return Now().UnixMilliStr()
}
func UnixMicroStr() string {
	return Now().UnixMicroStr()
}
func UnixNanoStr() string {
	return Now().UnixNanoStr()
}
func Date() string {
	return time.Now().Format(time.DateOnly)
}
func Datetime() string {
	return time.Now().Format(time.DateTime)
}

func TodayZeroMilli() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t.UnixNano() / 1e6
}

func Now() *Time {
	return &Time{
		time.Now(),
	}
}
func MustParseTime(v any) time.Time {
	t, err := parseTime(v)
	if err != nil { //如果传参解析失败，返回空时间
		return time.Time{}
	}
	return t
}
func ParseTime(v any) (time.Time, error) {
	t, err := parseTime(v)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// 还可以加个layout
func parseTime(v any) (time.Time, error) {
	switch t := v.(type) {
	case time.Time:
		return t, nil
	case string:
		return parseTimeFromString(t)
	case []byte:
		return parseTimeFromString(string(t))
	case int:
		return time.Unix(int64(t), 0), nil
	case int32:
		return time.Unix(int64(t), 0), nil
	case int64:
		return time.Unix(t, 0), nil
	default:
		return time.Time{}, errors.New("unsupported type")
	}
}
func parseTimeFromString(s string) (time.Time, error) {
	if ts, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(ts, 0), nil
	}
	loc := time.Now().Location()
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, loc); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("invalid time format")
}
