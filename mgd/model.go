package mgd

import (
	"context"
)

type CollectionNameGetter interface {
	CollectionName() string
}

type Model interface {
	GetId() interface{}
	SetId(id interface{})
	Create(id interface{})
	Update(id interface{})
}
type Page struct {
	PageNum  int64 `json:"pageNum"`
	PageSize int64 `json:"pageSize"`
	Total    int64 `json:"total"`
}
type PageList[T any] struct {
	Page Page `json:"page"`
	List []T  `json:"list"`
}
type InsertOptions struct {
	Context  context.Context
	Document interface{} //需要传入数组
}

type UpdateOptions struct {
	Context           context.Context
	Filter            interface{}
	Set               interface{}
	UnSet             interface{}
	AddToSet          interface{}
	Push              interface{}
	Inc               interface{}
	NoFoundError      bool
	ReturnNewDocument bool //返回新的文档
	Upset             bool
}

type DeleteOptions struct {
	Context context.Context
	Filter  interface{}
	One     bool
}
type Example struct {
	Context   context.Context
	Filter    interface{} //条件，只能是g.map(0值不忽略)或struct(0值忽略)
	Limit     int64       //限制条数
	Asc       []string    //正序
	Desc      []string    //倒叙
	Select    []string    //需要显示的字段
	Exclude   []string    //不需要显示的字段
	Page      Page
	BatchSize int32 //单批获取数据大小,当使用后台分发集群，使用连接池时每次getMore可能会分发到其他机器导致拿不到Cursor报错
}
