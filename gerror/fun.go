package gerror

import "fmt"

var (
	//系统前台可判断级不会动
	//-1通用msg，-2弹窗消息
	ServerError  = func() Error { return New(-10, "服务器繁忙").WithSkip(2) } //默认错误
	UnAuth       = func() Error { return New(-11, "未授权").WithSkip(2) }
	LoginTimeout = func() Error { return New(-12, "登录过期").WithSkip(2) }
	//错误码可能会动
	InvalidParam   = func() Error { return New(-20, "无效的参数").WithSkip(2) }
	BusinessError  = func() Error { return New(-21, "业务错误").WithSkip(2) }
	IllegalRequest = func() Error { return New(-22, "非法请求").WithSkip(2) }
	AuthDenied     = func() Error { return New(-23, "权限不足").WithSkip(2) }
	IllegalSign    = func() Error { return New(-24, "无效签名").WithSkip(2) }
)

func New(code int, msg string) Error {
	return new(Error).WithCode(code).WithMsg(msg)
}

// 通用消息
func Msg(msg string) Error {
	return new(Error).WithCode(-1).WithMsg(msg).WithSkip(1)
}
func NewErrorSkip1(e error) Error {
	return new(Error).WithError(e).WithSkip(2) //包装了一层本身已+1
}

// 前台需要弹窗确认的错误
func Confirm(msg string) Error {
	return new(Error).WithCode(-2).WithMsg(msg).WithSkip(1)
}
func (err Error) WithCode(code int) Error {
	err.Code = code
	return err
}
func (err Error) WithMsg(msg ...any) Error {
	err.Msg = fmt.Sprint(msg...)
	return err
}
func (err Error) WithDetail(msg ...any) Error {
	err.Detail = fmt.Sprint(msg...)
	return err
}
func (err Error) WithError(e error) Error {
	err.Cause = e
	return err
}
func (err Error) WithSkip(skip ...int) Error {
	if len(skip) > 0 {
		err.Stack = callers(skip[0])
	} else {
		err.Stack = callers()
	}
	return err
}
