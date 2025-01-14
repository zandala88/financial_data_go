package server

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

func CronDailyWorker() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("Error loading Shanghai timezone:", err)
		return
	}
	c := cron.New(cron.WithLocation(location))

	c.AddFunc("0 23 * * *", func() {
		zap.S().Debugf("[CronDailyWorker] [start]")

		//tushare.Daily(context.Background(), "")
	})

	c.Start()
	defer c.Stop()
	select {}
}
