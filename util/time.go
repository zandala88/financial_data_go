package util

import (
	"strings"
	"time"
)

// RFC3339ToNormalTime
// RFC3339 日期格式标准化
func RFC3339ToNormalTime(rfc3339 string) string {
	if len(rfc3339) < 19 || rfc3339 == "" || !strings.Contains(rfc3339, "T") {
		return rfc3339
	}
	return strings.Split(rfc3339, "T")[0] + " " + strings.Split(rfc3339, "T")[1][:8]
}

func ConvertDateStrToTime(dateStr string, layout string) time.Time {
	// 使用 time.Parse 方法转换
	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}
	}

	// 返回转换后的 time.Time 类型
	return parsedTime
}
