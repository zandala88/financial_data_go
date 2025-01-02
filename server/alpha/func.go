package alpha

import "time"

func getYesterdayStr() string {
	return time.Now().AddDate(0, 0, -1).Format("2006-01-02")
}

func getNextWeekday() time.Time {
	// 获取当前时间
	today := time.Now()

	// 计算明天的日期
	tomorrow := today.AddDate(0, 0, 1)

	// 检查明天是否是周六或者周日，如果是，则跳过
	for tomorrow.Weekday() == time.Saturday || tomorrow.Weekday() == time.Sunday {
		tomorrow = tomorrow.AddDate(0, 0, 1) // 跳到下一个工作日
	}

	// 返回明天日期的 Unix 时间戳（秒）
	return tomorrow
}
