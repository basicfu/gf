package ghash

import (
	"crypto/sha256"
	"encoding/hex"
)

// hash不属于加密，暂时放在crypto目录下
func Sha256(v string) string {
	bs := sha256.Sum256([]byte(v))
	return hex.EncodeToString(bs[:])
}
