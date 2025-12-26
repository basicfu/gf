package gmd5

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

type ByteString interface {
	~string | ~[]byte
}

func Encrypt[T ByteString](data T) (encrypt string, err error) {
	return EncryptBytes(toBytes(data))
}

func MustEncrypt[T ByteString](data T) string {
	result, err := Encrypt(data)
	if err != nil {
		panic(err)
	}
	return result
}

func EncryptBytes(data []byte) (encrypt string, err error) {
	h := md5.New()
	if _, err = h.Write([]byte(data)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func MustEncryptBytes(data []byte) string {
	result, err := EncryptBytes(data)
	if err != nil {
		panic(err)
	}
	return result
}

func EncryptString(data string) (encrypt string, err error) {
	return EncryptBytes([]byte(data))
}

func MustEncryptString(data string) string {
	result, err := EncryptString(data)
	if err != nil {
		panic(err)
	}
	return result
}

func toBytes[T ByteString](data T) []byte {
	switch v := any(data).(type) {
	case string:
		return []byte(v)
	case []byte:
		return v
	default:
		return nil
	}
}

func EncryptFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func MustEncryptFile(path string) string {
	result, err := EncryptFile(path)
	if err != nil {
		panic(err)
	}
	return result
}
