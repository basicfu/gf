package gutil

import (
	"github.com/basicfu/gf/guid"
)

type Map = map[string]interface{}

func Panic(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Rid  string      `json:"rid"`
}
type Page struct {
	Total    int64
	PageSize int64
	PageNum  int64
}
type RollPage struct {
	Total     int64
	PageSize  int64
	NextToken string
}

func Success(data interface{}) Result {
	return Result{
		Code: 0,
		Msg:  "",
		Data: data,
		Rid:  guid.New(),
	}
}
func Error(code int, msg string) Result {
	return Result{
		Code: code,
		Msg:  msg,
		Data: nil,
		Rid:  guid.New(),
	}
}
