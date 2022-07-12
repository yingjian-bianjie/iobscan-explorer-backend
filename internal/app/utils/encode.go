package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func Sha1(v, salt string) string {
	r := sha1.Sum([]byte(salt + v))
	return hex.EncodeToString(r[:])
}

func Md5(v string) string {
	r := md5.Sum([]byte(v))
	return hex.EncodeToString(r[:])
}
