package gbase64

import (
	"encoding/base64"
)

func Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func EncodeString(src string) string {
	return string(Encode([]byte(src)))
}
func Decode(data []byte) ([]byte, error) {
	src := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(src, data)
	return src[:n], err
}

func DecodeString(data string) (string, error) {
	d, err := Decode([]byte(data))
	return string(d), err
}

func DecodeStringM(data string) string {
	result, err := DecodeString(data)
	if err != nil {
		panic(err)
	}
	return result
}
