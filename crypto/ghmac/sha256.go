package ghmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func Sha256(v, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(v))
	return h.Sum(nil)
}
func Sha256Hex(v, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(v))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}
func Sha256Base64(v, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(v))
	sum := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}
func Sha256Base64Raw(v, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(v))
	sum := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(sum)
}
