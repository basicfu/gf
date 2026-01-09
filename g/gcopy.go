package g

import (
	"github.com/jinzhu/copier"
)

/************* struct<=>struct,结构体-结构体，或直接传入传出数组结构体 *************/
//待添加option里的deepcopy功能，视实际使用情况封装，可加到第三个...参数option对象直接
func Copy[T any](input any, pointer T) T {
	_ = copier.Copy(&pointer, input)
	return pointer
}
func CopyE[T any](input any, pointer T) (T, error) {
	err := copier.Copy(&pointer, input)
	if err != nil {
		return pointer, err
	}
	return pointer, nil
}

//结构体合并用 dario.cat/mergo
