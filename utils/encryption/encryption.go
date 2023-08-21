// Package encryption @Author: youngalone [2023/8/20]
package encryption

import (
	"crypto/md5"
	"encoding/hex"
)

func Encrypt(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
