package public

import "time"

func GetTodayStr() string {
	return time.Now().Format("2006-01-02")
}
