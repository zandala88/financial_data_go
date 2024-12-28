package server

import (
	"financia/server/alpha"
	"github.com/robfig/cron/v3"
)

func InsertDailyDate() {
	c := cron.New()

	// 添加每天早上5点执行的任务
	c.AddFunc("0 11 * * *", func() {
		alpha.AlphaDaily()
	})

	// 启动定时器
	c.Start()

	// 保持程序运行，避免退出
	select {}
}
