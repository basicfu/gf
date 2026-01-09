package gresult

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func New(code int, msg string, data any) Result {
	return Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
func Success(data any) Result {
	return Result{
		Code: 0,
		Msg:  "",
		Data: data,
	}
}
func Error(code int, msg string) Result {
	return Result{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
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
