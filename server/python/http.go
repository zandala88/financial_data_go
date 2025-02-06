package python

import (
	"context"
	"financia/config"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/model"
	"financia/util"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"math"
	"time"
)

func PythonPredictStock(id int, stockData []*model.StockData) (float64, error) {
	pyReq := &PythonPredictReq{
		Data: make([]*PythonPredictReqSimple, 0, len(stockData)),
	}

	for _, v := range stockData {
		pyReq.Data = append(pyReq.Data, &PythonPredictReqSimple{
			Date:   v.TradeDate.Format(time.DateOnly),
			CoIMF1: v.Open,
			CoIMF2: v.High,
			CoIMF3: v.Low,
			CoIMF4: v.Vol,
			Target: v.Close,
		})
	}

	pyResp := &PythonPredictResp{}

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(pyReq).
		SetResult(&pyResp).
		Post(config.Configs.Python.Url)

	if err != nil {
		zap.S().Error("[PredictStock] [err] = ", err.Error())
		return 0, err
	}

	pyResp.Data.Val = math.Floor(pyResp.Data.Val*1000) / 1000

	go func() {
		rdb := connector.GetRedis()
		rdb.Set(context.Background(), fmt.Sprintf(public.RedisKeyStockPredict, id), pyResp.Data.Val, time.Second*time.Duration(util.SecondsUntilMidnight()))
	}()

	return pyResp.Data.Val, nil
}

func PythonPredictFund(id int, fundData []*model.FundData) (float64, error) {
	pyReq := &PythonPredictReq{
		Data: make([]*PythonPredictReqSimple, 0, len(fundData)),
	}

	for _, v := range fundData {
		pyReq.Data = append(pyReq.Data, &PythonPredictReqSimple{
			Date:   v.TradeDate.Format(time.DateOnly),
			CoIMF1: v.Open,
			CoIMF2: v.High,
			CoIMF3: v.Low,
			CoIMF4: 0,
			Target: v.Close,
		})
	}

	pyResp := &PythonPredictResp{}

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(pyReq).
		SetResult(&pyResp).
		Post(config.Configs.Python.Url)

	if err != nil {
		zap.S().Error("[PredictStock] [err] = ", err.Error())
		return 0, err
	}

	pyResp.Data.Val = math.Floor(pyResp.Data.Val*1000) / 1000

	go func() {
		rdb := connector.GetRedis()
		rdb.Set(context.Background(), fmt.Sprintf(public.RedisKeyFundPredict, id), pyResp.Data.Val, time.Second*time.Duration(util.SecondsUntilMidnight()))
	}()

	return pyResp.Data.Val, nil
}
