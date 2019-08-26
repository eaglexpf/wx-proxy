package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func Md5(secret string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1(secret string) string {
	h := sha1.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(nil))
}
