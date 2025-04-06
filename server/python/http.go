package python

import (
	"context"
	"encoding/json"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/dao"
	"financia/public/db/model"
	pb "financia/server/python/grpc"
	"financia/server/tushare"
	"financia/util"
	"fmt"
	"go.uber.org/zap"
	"math"
	"sort"
	"time"
)

func PythonPredictStock(id int, stockData []*model.StockData) (float64, error) {
	pyReq := &pb.PredictRequest{
		Data: make([]*pb.DataPoint, 0, len(stockData)),
	}

	for _, v := range stockData {
		pyReq.Data = append(pyReq.Data, &pb.DataPoint{
			Date:   v.TradeDate.Format(time.DateOnly),
			CoImf1: v.Open,
			CoImf2: v.High,
			CoImf3: v.Low,
			CoImf4: float64(v.Vol),
			Target: v.Close,
		})
	}

	val, err := SendPredictRequest(pyReq)
	if err != nil {
		zap.S().Error("[PredictStock] [err] = ", err.Error())
		return 0, err
	}

	val = math.Floor(val*1000) / 1000

	go func() {
		rdb := connector.GetRedis()
		rdb.Set(context.Background(), fmt.Sprintf(public.RedisKeyStockPredict, id), val, time.Second*time.Duration(util.SecondsUntilMidnight()))

		tsCode := stockData[0].TsCode

		// 找到下一个交易日
		type CalFutResp struct {
			Sse  []*tushare.FutTradeCalResp `json:"sse"`
			Szse []*tushare.FutTradeCalResp `json:"szse"`
		}
		resp := &CalFutResp{
			Sse:  make([]*tushare.FutTradeCalResp, 0),
			Szse: make([]*tushare.FutTradeCalResp, 0),
		}
		result, _ := rdb.Get(context.Background(), "cal_fut").Result()
		json.Unmarshal([]byte(result), resp)
		sort.Slice(resp.Sse, func(i, j int) bool {
			return resp.Sse[i].CalDate < resp.Sse[j].CalDate
		})
		var nextTradeDate string
		for _, v := range resp.Sse {
			if v.IsOpen == 1 && util.ConvertDateStrToTime(v.CalDate, time.DateOnly).Before(stockData[len(stockData)-1].TradeDate) {
				zap.S().Debugf("[PredictStock] [nextTradeDate] = %v, %v", v.CalDate, stockData[len(stockData)-1].TradeDate)
				nextTradeDate = v.CalDate
				break
			}
		}

		dao.InsertStockPredict(context.Background(), &model.StockPredict{
			TsCode:    tsCode,
			TradeDate: util.ConvertDateStrToTime(nextTradeDate, time.DateOnly),
			Predict:   val,
		})
	}()

	return val, nil
}

func PythonPredictFund(id int, fundData []*model.FundData) (float64, error) {
	pyReq := &pb.PredictRequest{
		Data: make([]*pb.DataPoint, 0, len(fundData)),
	}

	for _, v := range fundData {
		pyReq.Data = append(pyReq.Data, &pb.DataPoint{
			Date:   v.TradeDate.Format(time.DateOnly),
			CoImf1: v.Open,
			CoImf2: v.High,
			CoImf3: v.Low,
			CoImf4: 0,
			Target: v.Close,
		})
	}

	val, err := SendPredictRequest(pyReq)
	if err != nil {
		zap.S().Error("[PredictStock] [err] = ", err.Error())
		return 0, err
	}

	val = math.Floor(val*1000) / 1000

	go func() {
		rdb := connector.GetRedis()
		rdb.Set(context.Background(), fmt.Sprintf(public.RedisKeyFundPredict, id), val, time.Second*time.Duration(util.SecondsUntilMidnight()))

		tsCode := fundData[0].TsCode

		// 找到下一个交易日
		type CalFutResp struct {
			Sse  []*tushare.FutTradeCalResp `json:"sse"`
			Szse []*tushare.FutTradeCalResp `json:"szse"`
		}
		resp := &CalFutResp{
			Sse:  make([]*tushare.FutTradeCalResp, 0),
			Szse: make([]*tushare.FutTradeCalResp, 0),
		}
		result, _ := rdb.Get(context.Background(), "cal_fut").Result()
		json.Unmarshal([]byte(result), resp)
		sort.Slice(resp.Sse, func(i, j int) bool {
			return resp.Sse[i].CalDate < resp.Sse[j].CalDate
		})
		var nextTradeDate string
		for _, v := range resp.Sse {
			if v.IsOpen == 1 && util.ConvertDateStrToTime(v.CalDate, time.DateOnly).Before(fundData[len(fundData)-1].TradeDate) {
				nextTradeDate = v.CalDate
				break
			}
		}

		dao.InsertFundPredict(context.Background(), &model.FundPredict{
			TsCode:    tsCode,
			TradeDate: util.ConvertDateStrToTime(nextTradeDate, time.DateOnly),
			Predict:   val,
		})
	}()

	return val, nil
}
