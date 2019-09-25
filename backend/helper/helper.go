package helper

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

func NowSecond() int64 {
	return time.Now().Unix()
}

func NowMs() int64 {
	return time.Now().UnixNano() / 1000000
}

func GenerateRandomBytes(size int) ([]byte, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %v", err)
	}
	return buf, nil
}

func Base64(bs []byte) string {
	return base64.StdEncoding.EncodeToString(bs)
}

func Base64Decode(s string) (x []byte, err error) {
	x, err = base64.StdEncoding.DecodeString(s)
	return
}

func MD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
