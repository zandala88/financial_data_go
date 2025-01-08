package public

import (
	"math/rand"
	"time"
)

func RandomWithinPercentage(value float64, percentage float64) float64 {
	// 使用随机数源创建随机数生成器
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 计算波动范围
	offset := value * percentage / 100
	// 生成波动范围内的随机浮点数
	return value + (source.Float64()*2*offset - offset)
}
