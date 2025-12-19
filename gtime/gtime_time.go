package gtime

import (
	"strconv"
	"time"
)

var layouts = []string{
	time.Layout,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,
	// 自定义常用格式
	"2006/01/02 15:04:05",
	"2006/01/02",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05.000Z07:00",
	"2006-01-02 15:04:05 -0700",
}

type Time struct {
	time.Time
}

func (t Time) String() string {
	if t.IsZero() {
		return ""
	}
	if t.Year() == 0 {
		return t.Format("15:04:05")
	}
	return t.Format("2006-01-02 15:04:05")
}
func New(v ...any) Time {
	if len(v) == 0 {
		return Time{
			time.Now(),
		}
	}
	t, err := parseTime(v[0])
	if err != nil { //如果传参解析失败，返回空时间
		return Time{}
	}
	return Time{t}
}

func (t Time) UnixStr() string {
	if t.IsZero() {
		return ""
	}
	return strconv.FormatInt(t.Unix(), 10)
}

func (t Time) UnixMilliStr() string {
	if t.IsZero() {
		return ""
	}
	return strconv.FormatInt(t.UnixMilli(), 10)
}
func (t Time) UnixMicroStr() string {
	if t.IsZero() {
		return ""
	}
	return strconv.FormatInt(t.UnixMicro(), 10)
}
func (t Time) UnixNanoStr() string {
	if t.IsZero() {
		return ""
	}
	return strconv.FormatInt(t.UnixNano(), 10)
}
func (t Time) Date() string {
	return t.Time.Format(time.DateOnly)
}
func (t Time) DateTime() string {
	return t.Time.Format(time.DateTime)
}

func (t Time) Add(d time.Duration) Time {
	return Time{t.Time.Add(d)}
}
func (t Time) AddDate(years int, months int, days int) Time {
	return Time{t.Time.AddDate(years, months, days)}
}
func (t Time) After(u Time) bool {
	return t.Time.After(u.Time)
}
func (t Time) Before(u Time) bool {
	return t.Time.Before(u.Time)
}
func (t Time) Equal(u Time) bool {
	return t.Time.Equal(u.Time)
}
