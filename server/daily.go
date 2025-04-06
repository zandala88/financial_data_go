package server

import (
	"context"
	"encoding/json"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/dao"
	"financia/public/db/model"
	"financia/server/python"
	"financia/server/tushare"
	"financia/service/fut"
	"financia/util"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sort"
	"strings"
	"sync"
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

	c.AddFunc("0 8 * * *", DailyPredictBefore)
	c.AddFunc("0 10 * * *", DailyPredict)

	c.Start()
	defer c.Stop()
	select {}
}

func DailyPredictBefore() {
	ctx := context.Background()
	rdb := connector.GetRedis().WithContext(ctx)
	tsCodeList, _ := rdb.SMembers(ctx, public.RedisKeyPredictList).Result()
	wg := sync.WaitGroup{}
	for _, tsCode := range tsCodeList {
		stockData, _ := dao.GetStockDataLimit30(ctx, tsCode)
		last := stockData[len(stockData)-1]
		date := strings.ReplaceAll(last.TradeDate.Add(time.Hour*24).Format(time.DateOnly), "-", "")
		data := tushare.DailyStockAll(ctx, &tushare.DailyReq{
			TsCode:    tsCode,
			StartDate: date,
		})
		_ = dao.InsertStockData(ctx, data)

		// 异步比较今日已更新
		wg.Add(1)
		go func(tsCode string) {
			defer wg.Done()
			key := fmt.Sprintf(public.RedisKeyStockDataDoToday, tsCode)
			rdb.Set(ctx, key, "1", time.Duration(util.SecondsUntilMidnight())*time.Second)
		}(tsCode)
	}

	wg.Wait()
}

func DailyPredict() {
	ctx := context.Background()
	rdb := connector.GetRedis().WithContext(ctx)
	tsCodeList, _ := rdb.SMembers(ctx, public.RedisKeyPredictList).Result()
	db := connector.GetDB().WithContext(ctx)
	for _, tsCode := range tsCodeList {
		var id int
		db.Model(model.StockInfo{}).Select("f_id").Where("f_ts_code = ?", tsCode).Scan(&id)
		stockData, _ := dao.GetStockDataLimit30(ctx, tsCode)
		sort.Slice(stockData, func(i, j int) bool {
			return stockData[i].TradeDate.Before(stockData[j].TradeDate)
		})
		_, _ = python.PythonPredictStock(id, stockData)
	}
}
