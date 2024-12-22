package alpha

import "time"

func getYesterdayStr() string {
	return time.Now().AddDate(0, 0, -1).Format("2006-01-02")
}
