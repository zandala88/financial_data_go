package server

import (
	"context"
	"encoding/json"
	"financia/public/db/connector"
	"financia/server/tushare"
	"financia/service/fut"
	"financia/util"
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

	c.AddFunc("0 0 * * *", func() {

		ctx := context.Background()

		resp := &fut.CalFutResp{
			Sse:  make([]*tushare.FutTradeCalResp, 0),
			Szse: make([]*tushare.FutTradeCalResp, 0),
		}
		resp.Sse, resp.Szse = tushare.FutTradeCal(ctx)

		for i, v := range resp.Sse {
			resp.Sse[i].CalDate = util.ConvertDateStrToTime(v.CalDate, util.TimeDateOnlyWithOutSep).Format(time.DateOnly)
		}
		for i, v := range resp.Szse {
			resp.Szse[i].CalDate = util.ConvertDateStrToTime(v.CalDate, util.TimeDateOnlyWithOutSep).Format(time.DateOnly)
		}
		// 当天0点过期
		exp := util.SecondsUntilMidnight()
		rdbStr, _ := json.Marshal(resp)

		rdb := connector.GetRedis().WithContext(ctx)
		_, err := rdb.Set(ctx, "cal_fut", rdbStr, time.Duration(exp)*time.Second).Result()
		if err != nil {
			zap.S().Errorf("[CalFut] [rdb.Set] [err] = %s", err.Error())
		}
	})

	c.Start()
	defer c.Stop()
	select {}
}
