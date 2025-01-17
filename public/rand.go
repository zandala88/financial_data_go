package public

import (
	"crypto/md5"
	"encoding/hex"
	"math"
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

func GetRandomSliceWithSqrtLength(len int) (start, end int) {
	// 计算数组长度的平方根，并向下取整
	sqrtLength := int(math.Sqrt(float64(len)))

	// 设置随机数种子
	rand.NewSource(time.Now().UnixNano())

	// 生成一个随机的起始索引，确保起始索引加上平方根长度不超过数组的长度
	startIndex := rand.Intn(len - sqrtLength + 1)

	return startIndex, startIndex + sqrtLength
}
