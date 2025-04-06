package python

import (
	"context"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/model"
	pb "financia/server/python/grpc"
	"financia/util"
	"fmt"
	"go.uber.org/zap"
	"math"
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
	}()

	return val, nil
}

func PythonPredictAllStock(_ int, stockData []*model.StockData) ([]float64, error) {
	pyReq := &pb.PredictAllRequest{
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

	list, err := SendPredictAllRequest(pyReq)
	if err != nil {
		zap.S().Error("[PredictStock] [err] = ", err.Error())
		return nil, err
	}
	return list, nil
}

func PythonPredictFund(_ int, fundData []*model.FundData) (float64, error) {
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

	return val, nil
}
