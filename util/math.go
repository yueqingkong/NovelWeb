package util

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5
func MD5(value string) string {
	hash := md5.New()
	hash.Write([]byte(value))
	return hex.EncodeToString(hash.Sum(nil))
}
