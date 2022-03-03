package md5

import (
	"crypto/md5"
	"encoding/hex"
)

// Encode md5 encryption
func Encode(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
