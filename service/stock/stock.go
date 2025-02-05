package stock

import (
	"encoding/json"
	"errors"
	"financia/config"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/dao"
	"financia/server/tushare"
	"financia/service/fut"
	"financia/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"math"
	"sort"
	"strings"
	"time"
)

func DataStock(c *gin.Context) {
	var req DataStockReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[DataStock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	zap.S().Debugf("[DataStock] [req] = %#v", req)

	info, err := dao.GetStockInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[DataStock] [GetStockInfo] [err] = ", err.Error())
		return
	}

	list, err := dao.GetStockData(c, info.TsCode, req.StartDate, req.EndDate)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[DataStock] [GetStockData] [err] = ", err.Error())
		return
	}

	if len(list) == 0 {
		data := tushare.DailyStockAll(c, &tushare.DailyReq{
			TsCode: info.TsCode,
		})
		if err := dao.InsertStockData(c, data); err != nil {
			zap.S().Error("[DataStock] [InsertStockData] [err] = ", err.Error())
		}
		list, err = dao.GetStockData(c, info.TsCode, req.StartDate, req.EndDate)
	}

	respList := make([]*DataStockSimple, 0, len(list))
	for _, v := range list {
		respList = append(respList, &DataStockSimple{
			TradeDate: v.TradeDate.Format(time.DateOnly),
			Open:      v.Open,
			High:      v.High,
			Low:       v.Low,
			Close:     v.Close,
			PreClose:  v.PreClose,
			Change:    v.Change,
			PctChg:    v.PctChg,
			Vol:       v.Vol,
			Amount:    v.Amount,
		})
	}

	// 异步更新数据
	go func() {
		last := list[len(list)-1]
		date := strings.ReplaceAll(last.TradeDate.Add(time.Hour*24).Format(time.DateOnly), "-", "")
		data := tushare.DailyStockAll(c, &tushare.DailyReq{
			TsCode:    info.TsCode,
			StartDate: date,
		})
		if err := dao.InsertStockData(c, data); err != nil {
			zap.S().Error("[DataStock] [InsertStockData] [err] = ", err.Error())
		}
	}()

	util.SuccessResp(c, &DataStockResp{
		Have: true,
		List: respList,
	})

}

func GraphStock(c *gin.Context) {
	rdb := connector.GetRedis().WithContext(c)
	result, err := rdb.Get(c, public.RedisKeyGraphStock).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GraphStock] [rdb.Get] [err] = ", err.Error())
		return
	}

	if errors.Is(err, redis.Nil) {
		fields, err := dao.CountStockFields(c)
		if err != nil {
			util.FailRespWithCode(c, util.InternalServerError)
			zap.S().Errorf("CountStockFields error: %s", err.Error())
			return
		}

		resp := &GraphStockResp{
			IsHs:     fields["is_hs"],
			Exchange: fields["exchange"],
			Market:   fields["market"],
		}

		go func() {
			listStr, _ := json.Marshal(resp)
			_, err := rdb.Set(c, public.RedisKeyGraphStock, listStr, 0).Result()
			if err != nil {
				util.FailRespWithCode(c, util.InternalServerError)
				zap.S().Error("[GraphStock] [rdb.Set] [err] = ", err.Error())
				return
			}
		}()

		util.SuccessResp(c, resp)
		return
	}

	resp := &GraphStockResp{}
	if err := json.Unmarshal([]byte(result), resp); err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GraphStock] [json.Unmarshal] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, &resp)
}

func HaveStock(c *gin.Context) {
	var req HaveStockReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[DataStock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	zap.S().Debugf("[DataStock] [req] = %#v", req)

	info, err := dao.GetStockInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[DataStock] [GetStockInfo] [err] = ", err.Error())
		return
	}

	have, err := dao.CheckStockData(c, info.TsCode)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[DataStock] [GetStockData] [err] = ", err.Error())
		return
	}

	if !have {
		data := tushare.DailyStockAll(c, &tushare.DailyReq{
			TsCode: info.TsCode,
		})
		have = len(data) > 0
		if err := dao.InsertStockData(c, data); err != nil {
			util.FailRespWithCode(c, util.InternalServerError)
			zap.S().Errorf("[Daily] [InsertStockData] [err] = %s", err.Error())
			return
		}
	}

	util.SuccessResp(c, &HaveStockResp{
		Have: have,
	})
	return
}

func IncomeStock(c *gin.Context) {
	var req IncomeStockReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[IncomeStock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	stockInfo, err := dao.GetStockInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[IncomeStock] [GetStockInfo] [err] = ", err.Error())
		return
	}

	incomeList := tushare.StockIncome(c, stockInfo.TsCode)
	if incomeList == nil || len(incomeList) == 0 {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[IncomeStock] [StockIncome] [err] = ", "incomeList is nil")
		return
	}

	util.SuccessResp(c, &IncomeStockResp{
		List: incomeList,
	})
}

func InfoStock(c *gin.Context) {
	var req InfoStockReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[InfoStock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	info, err := dao.GetStockInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[InfoStock] [GetStockInfo] [err] = ", err.Error())
		return
	}

	rdb := connector.GetRedis().WithContext(c)
	userId := util.GetUid(c)
	redisKey := fmt.Sprintf(public.RedisKeyStockFollow, userId)
	zap.S().Debugf("[InfoStock] [redisKey] = %s", redisKey)
	follow, err := rdb.SIsMember(c, redisKey, req.Id).Result()
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[InfoStock] [rdb.SIsMember] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, &InfoStockResp{
		FullName: info.FullName,
		Industry: info.Industry,
		Market:   info.Market,
		Follow:   follow,
	})
}

func ListStock(c *gin.Context) {
	var req ListStockReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[ListStock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	zap.S().Debugf("[ListStock] [req] = %#v", req)

	list, count, err := dao.GetStockList(c, req.Search, req.IsHs, req.Exchange, req.Market, req.Page, req.PageSize)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[ListStock] [GetStockList] [err] = ", err.Error())
		return
	}

	respList := make([]*ListStockSimple, 0, len(list))
	for _, v := range list {
		respList = append(respList, &ListStockSimple{
			Id:         v.Id,
			Name:       v.Name,
			Area:       v.Area,
			Industry:   v.Industry,
			Market:     v.Market,
			ActName:    v.ActName,
			ActEntType: v.ActEntType,
			FullName:   v.FullName,
			EnName:     v.EnName,
			CnSpell:    v.CnSpell,
			Exchange:   v.Exchange,
			CurrType:   v.CurrType,
			ListStatus: v.ListStatus,
			IsHs:       v.IsHs,
		})
	}

	util.SuccessResp(c, &ListStockResp{
		List:         respList,
		HasMore:      count > int64(req.Page*(req.PageSize-1)+len(list)),
		TotalPageNum: int(count/int64(req.PageSize) + 1),
	})
}

func QueryStock(c *gin.Context) {
	fields, err := dao.DistinctStockFields(c)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Errorf("DistinctStockFields error: %s", err.Error())
		return
	}

	util.SuccessResp(c, &QueryStockResp{
		IsHsList:     fields["is_hs"],
		ExchangeList: fields["exchange"],
		MarketList:   fields["market"],
	})
}

func ForecastStock(c *gin.Context) {
	var req ForecastStockReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[ForecastStock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	stockInfo, err := dao.GetStockInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[ForecastStock] [GetStockInfo] [err] = ", err.Error())
		return
	}

	forecast := tushare.StockForecast(c, stockInfo.TsCode)
	if forecast == nil || len(forecast) == 0 {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[ForecastStock] [StockForecast] [err] = ", "forecast is nil")
		return
	}

	sort.Slice(forecast, func(i, j int) bool {
		return forecast[i].AnnDate < forecast[j].AnnDate
	})

	util.SuccessResp(c, &ForecastStockResp{
		List: forecast,
	})
}

func Top10Stock(c *gin.Context) {
	var req Top10StockReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[Top10Stock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	stockInfo, err := dao.GetStockInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[Top10Stock] [GetStockInfo] [err] = ", err.Error())
		return
	}

	top10 := tushare.StockHolderTop10(c, stockInfo.TsCode)
	if top10 == nil || len(top10) == 0 {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[Top10Stock] [StockHolderTop10] [err] = ", "top10 is nil")
		return
	}

	topDate := top10[0].AnnDate
	rank := make([]*Top10StockRank, 0)
	for _, v := range top10 {
		if v.AnnDate != topDate {
			break
		}
		rank = append(rank, &Top10StockRank{
			HoldRatio:  v.HoldRatio,
			HolderName: v.HolderName,
		})
	}

	sort.Slice(top10, func(i, j int) bool {
		return top10[i].AnnDate < top10[j].AnnDate
	})

	util.SuccessResp(c, &Top10StockResp{
		Rank: rank,
		List: top10,
	})
}

func Top10HsgtStock(c *gin.Context) {
	// 获取最近的交易日
	rdb := connector.GetRedis().WithContext(c)
	result, _ := rdb.Get(c, "cal_fut").Result()

	timeList := &fut.CalFutResp{
		Sse:  make([]*tushare.FutTradeCalResp, 0),
		Szse: make([]*tushare.FutTradeCalResp, 0),
	}
	json.Unmarshal([]byte(result), timeList)

	date := time.Now().Format(util.TimeDateOnlyWithOutSep)
	now := time.Now().Add(time.Hour * -24)
	for _, v := range timeList.Sse {
		t := util.ConvertDateStrToTime(v.CalDate, time.DateOnly)
		if t.After(now) {
			continue
		}

		if v.IsOpen == public.MarketStatusOpen {
			date = t.Format(util.TimeDateOnlyWithOutSep)
			break
		}
	}

	sh, sz := tushare.StockHsgtTop10(c, date)

	util.SuccessResp(c, &Top10HsgtStockResp{
		ShList: sh,
		SzList: sz,
	})
}

func PredictStock(c *gin.Context) {
	var req PredictStockReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[PredictStock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	stockInfo, err := dao.GetStockInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[PredictStock] [GetStockInfo] [err] = ", err.Error())
		return
	}

	stockData, err := dao.GetStockDataLimit30(c, stockInfo.TsCode)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[PredictStock] [GetStockDataLimit30] [err] = ", err.Error())
		return
	}

	if len(stockData) == 0 {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[PredictStock] [GetStockDataLimit30] [err] = ", "stockData is nil")
		return
	}

	sort.Slice(stockData, func(i, j int) bool {
		return stockData[i].TradeDate.Before(stockData[j].TradeDate)
	})

	last7 := make([]float64, 0, 7)
	for i := range stockData[len(stockData)-7:] {
		last7 = append(last7, stockData[i].Close)
	}

	rdb := connector.GetRedis().WithContext(c)
	result, err := rdb.Get(c, fmt.Sprintf(public.RedisKeyStockPredict, req.Id)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[PredictStock] [rdb.Get] [err] = ", err.Error())
		return
	}
	if !errors.Is(err, redis.Nil) {
		util.SuccessResp(c, &PredictStockResp{
			List: last7,
			Val:  cast.ToFloat64(result),
		})
		return
	}

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
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(pyReq).
		SetResult(&pyResp).
		Post(config.Configs.Python.Url)

	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[PredictStock] [err] = ", err.Error())
		return
	}

	pyResp.Data.Val = math.Floor(pyResp.Data.Val*1000) / 1000

	rdb.Set(c, fmt.Sprintf(public.RedisKeyStockPredict, req.Id), pyResp.Data.Val, time.Second*time.Duration(util.SecondsUntilMidnight()))

	util.SuccessResp(c, &PredictStockResp{
		List: last7,
		Val:  pyResp.Data.Val,
	})
}

func FollowStock(c *gin.Context) {
	var req FollowStockReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[FollowStock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	userId := util.GetUid(c)
	rdb := connector.GetRedis().WithContext(c)
	redisKey := fmt.Sprintf(public.RedisKeyStockFollow, userId)

	exists, err := rdb.SIsMember(c, redisKey, req.Id).Result()
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[FollowStock] [rdb.SIsMember] [err] = ", err.Error())
		return
	}
	if req.Follow == exists {
		util.SuccessResp(c, nil)
		return
	}

	if req.Follow {
		_, err = rdb.SAdd(c, redisKey, req.Id).Result()
		zap.S().Debugf("[FollowStock] [rdb.SAdd] [err] = %v", err)
		if err != nil {
			util.FailRespWithCode(c, util.InternalServerError)
			zap.S().Error("[FollowStock] [rdb.SAdd] [err] = ", err.Error())
			return
		}
	} else {
		_, err = rdb.SRem(c, redisKey, req.Id).Result()
		if err != nil {
			util.FailRespWithCode(c, util.InternalServerError)
			zap.S().Error("[FollowStock] [rdb.SRem] [err] = ", err.Error())
			return
		}
	}

	util.SuccessResp(c, nil)
}
