package gaes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrInvalidKey = errors.New("gaes: invalid key length (must be 16/24/32 bytes)")
	ErrInvalidMsg = errors.New("gaes: invalid ciphertext")
)

func EncryptM(key, plaintext string, add ...string) string {
	s, err := Encrypt(key, plaintext, add...)
	if err != nil {
		panic(err)
	}
	return s
}

func DecryptM(key, cipherText string, add ...string) string {
	b, err := Decrypt(key, cipherText, add...)
	if err != nil {
		panic(err)
	}
	return b
}

// aad常见场景如绑定用户，也可以传nil
func Encrypt(key, plaintext string, add ...string) (string, error) {
	keyB := []byte(key)
	var addB []byte
	if len(add) != 0 {
		addB = []byte(add[0])
	}
	if l := len(keyB); l != 16 && l != 24 && l != 32 {
		return "", ErrInvalidKey
	}
	block, err := aes.NewCipher(keyB)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), addB)
	out := make([]byte, 0, len(nonce)+len(ciphertext))
	out = append(out, nonce...)
	out = append(out, ciphertext...)
	return base64.StdEncoding.EncodeToString(out), nil
}

func Decrypt(key, cipherText string, add ...string) (string, error) {
	keyB := []byte(key)
	var addB []byte
	if len(add) != 0 {
		addB = []byte(add[0])
	}
	if l := len(keyB); l != 16 && l != 24 && l != 32 {
		return "", ErrInvalidKey
	}
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", ErrInvalidMsg
	}
	block, err := aes.NewCipher(keyB)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(data) < ns+1 {
		return "", ErrInvalidMsg
	}
	nonce := data[:ns]
	ciphertext := data[ns:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, addB)
	if err != nil {
		return "", ErrInvalidMsg
	}
	return string(plaintext), nil
}
