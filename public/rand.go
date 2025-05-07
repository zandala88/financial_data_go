package public

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

const charset = "0123456789"

func GenerateVerificationCode(length int) string {
	rand.NewSource(time.Now().UnixNano())
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func GenerateMD5Hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
