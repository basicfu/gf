package grsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/basicfu/gf/encoding/gbase64"
)

func Encrypt(publicKey, originalData string) string {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		panic("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	pub := pubInterface.(*rsa.PublicKey)
	v15, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(originalData))
	if err != nil {
		panic(err)
	}
	return string(gbase64.Encode(v15))
}

func Decrypt(privateKey, ciphertext string) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("private key error")
	}
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	data, err := gbase64.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	v15, err := rsa.DecryptPKCS1v15(nil, private, []byte(data))
	if err != nil {
		return "", err
	}
	return string(v15), nil
}

func MustDecrypt(privateKey, ciphertext string) string {
	result, err := Decrypt(privateKey, ciphertext)
	if err != nil {
		panic(err)
	}
	return result
}
