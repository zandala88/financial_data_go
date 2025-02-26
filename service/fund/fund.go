package fund

import (
	"encoding/json"
	"errors"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/dao"
	"financia/public/db/model"
	"financia/server/python"
	"financia/server/spark"
	"financia/server/tushare"
	"financia/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"sort"
	"strings"
	"time"
)

func DataFund(c *gin.Context) {
	var req DataFundReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[DataFund] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	info, err := dao.GetFundInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[DataFund] [GetFundInfo] [err] = ", err.Error())
		return
	}

	list, err := dao.GetFundData(c, info.TsCode, req.StartDate, req.EndDate)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[DataFund] [GetFundData] [err] = ", err.Error())
		return
	}

	if len(list) == 0 {
		data := tushare.DailyFundAll(c, &tushare.DailyReq{
			TsCode: info.TsCode,
		})
		if err := dao.InsertFundData(c, data); err != nil {
			zap.S().Error("[DataFund] [InsertStockData] [err] = ", err.Error())
		}
		list, err = dao.GetFundData(c, info.TsCode, req.StartDate, req.EndDate)
	}

	respList := make([]*DataFundSimple, 0, len(list))
	for _, v := range list {
		respList = append(respList, &DataFundSimple{
			TradeDate: v.TradeDate.Format(time.DateOnly),
			Open:      v.Open,
			High:      v.High,
			Low:       v.Low,
			Close:     v.Close,
			PreClose:  v.PreClose,
			Change:    v.Change,
			PctChg:    v.PctChg,
			Vol:       cast.ToInt64(v.Vol),
			Amount:    v.Amount,
		})
	}

	// 异步更新数据
	go func() {
		last := list[len(list)-1]
		date := strings.ReplaceAll(last.TradeDate.Add(time.Hour*24).Format(time.DateOnly), "-", "")
		data := tushare.DailyFundAll(c, &tushare.DailyReq{
			TsCode:    info.TsCode,
			StartDate: date,
		})
		if err := dao.InsertFundData(c, data); err != nil {
			zap.S().Error("[DataFund] [InsertFundData] [err] = ", err.Error())
		}
	}()

	rdb := connector.GetRedis().WithContext(c)
	userId := util.GetUid(c)
	redisKey := fmt.Sprintf(public.RedisKeyFundFollow, userId)
	follow, err := rdb.SIsMember(c, redisKey, req.Id).Result()
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[DataFund] [rdb.SIsMember] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, &DataFundResp{
		Follow: follow,
		Have:   true,
		List:   respList,
	})
}

func GraphFund(c *gin.Context) {
	resp := &GraphFundResp{}

	// 获取redis数据
	rdb := connector.GetRedis().WithContext(c)
	radioResult, err := rdb.Get(c, public.RedisKeyFundRadio).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[GraphFund] [rdb.Get] [err] = ", err.Error())
		return
	}
	if errors.Is(err, redis.Nil) {
		radio := tushare.FundSalesRatio(c)

		go func() {
			listStr, _ := json.Marshal(radio)
			if _, err := rdb.Set(c, public.RedisKeyFundRadio, listStr, time.Duration(util.SecondsUntilMidnight())*time.Second).Result(); err != nil {
				zap.S().Error("[GraphFund] [rdb.Set] [err] = ", err.Error())
				return
			}
		}()

		resp.Radio = radio
	} else {
		var radio []*tushare.FundSalesRatioResp
		if err := json.Unmarshal([]byte(radioResult), &radio); err != nil {
			util.FailRespWithCodeAndZap(c, util.InternalServerError, "[GraphFund] [json.Unmarshal] [err] = ", err.Error())
			return
		}
		resp.Radio = radio
	}

	volResult, err := rdb.Get(c, public.RedisKeyFundVol).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[GraphFund] [rdb.Get] [err] = ", err.Error())
		return
	}
	if errors.Is(err, redis.Nil) {
		vol := tushare.FundSalesVol(c)

		sort.Slice(vol, func(i, j int) bool {
			// year sec
			if vol[i].Year == vol[j].Year {
				// quarter sec
				if vol[i].Quarter == vol[j].Quarter {
					return vol[i].Rank < vol[j].Rank
				}
			}
			return vol[i].Year < vol[j].Year
		})

		go func() {
			listStr, _ := json.Marshal(vol)
			if _, err := rdb.Set(c, public.RedisKeyFundVol, listStr, time.Duration(util.SecondsUntilMidnight())*time.Second).Result(); err != nil {
				zap.S().Error("[GraphFund] [rdb.Set] [err] = ", err.Error())
				return
			}
		}()

		resp.Inst = vol
	} else {
		var vol []*tushare.FundSalesVolResp
		if err := json.Unmarshal([]byte(volResult), &vol); err != nil {
			util.FailRespWithCodeAndZap(c, util.InternalServerError, "[GraphFund] [json.Unmarshal] [err] = ", err.Error())
			return
		}
		resp.Inst = vol
	}

	util.SuccessResp(c, resp)
}

func HaveFund(c *gin.Context) {
	var req HaveFundReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[DataFund] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	info, err := dao.GetFundInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[DataFund] [GetFundInfo] [err] = ", err.Error())
		return
	}

	have, err := dao.CheckFundData(c, info.TsCode)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[DataFund] [CheckFundData] [err] = ", err.Error())
		return
	}

	if !have {
		data := tushare.DailyFundAll(c, &tushare.DailyReq{
			TsCode: info.TsCode,
		})
		have = len(data) > 0
		if err := dao.InsertFundData(c, data); err != nil {
			util.FailRespWithCodeAndZap(c, util.InternalServerError, "[DataFund] [InsertFundData] [err] = ", err.Error())
			return
		}
		if have {
			dao.UpdateFund(c, &model.FundInfo{
				Id:   int64(req.Id),
				Flag: 1,
			})
		}

	}

	util.SuccessResp(c, &HaveFundResp{
		Have: have,
	})
}

func ListFund(c *gin.Context) {
	var req ListFundReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[ListStock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	list, count, err := dao.GetFundList(c, req.Search, req.FundType, req.InvestType, req.Page, req.PageSize)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[ListStock] [GetStockList] [err] = ", err.Error())
		return
	}

	respList := make([]*ListFundSimple, 0, len(list))
	for _, v := range list {
		respList = append(respList, &ListFundSimple{
			Id:           v.Id,
			Name:         v.Name,
			Management:   v.Management,
			Custodian:    v.Custodian,
			FundType:     v.FundType,
			IssueAmount:  v.IssueAmount,
			MFree:        v.MFee,
			CFree:        v.CCFee,
			DurationYear: v.DurationYear,
			PValue:       v.PValue,
			MinAmount:    v.MinAmount,
			ExpReturn:    v.ExpReturn,
			Benchmark:    v.Benchmark,
			InvestType:   v.InvestType,
			Type:         v.Type,
			Trustee:      v.Trustee,
		})
	}

	util.SuccessResp(c, &ListFundResp{
		List:         respList,
		HasMore:      count > int64(req.Page*(req.PageSize-1)+len(list)),
		TotalPageNum: int(count/int64(req.PageSize) + 1),
	})
}

func QueryFund(c *gin.Context) {
	fields, err := dao.DistinctFundFields(c)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[QueryFund] [DistinctFundFields] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, &QueryFundResp{
		FundTypeList:   fields["fund_type"],
		InvestTypeList: fields["invest_type"],
	})
}

func FollowFund(c *gin.Context) {
	var req FollowFundReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[FollowFund] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	userId := util.GetUid(c)
	rdb := connector.GetRedis().WithContext(c)
	redisKey := fmt.Sprintf(public.RedisKeyFundFollow, userId)

	exists, err := rdb.SIsMember(c, redisKey, req.Id).Result()
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[FollowFund] [rdb.SIsMember] [err] = ", err.Error())
		return
	}
	if req.Follow == exists {
		util.SuccessResp(c, nil)
		return
	}

	if req.Follow {
		_, err = rdb.SAdd(c, redisKey, req.Id).Result()
		if err != nil {
			util.FailRespWithCodeAndZap(c, util.InternalServerError, "[FollowFund] [rdb.SAdd] [err] = ", err.Error())
			return
		}
	} else {
		_, err = rdb.SRem(c, redisKey, req.Id).Result()
		if err != nil {
			util.FailRespWithCodeAndZap(c, util.InternalServerError, "[FollowFund] [rdb.SRem] [err] = ", err.Error())
			return
		}
	}

	util.SuccessResp(c, nil)
}

func PredictFund(c *gin.Context) {
	var req PredictFundReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[PredictFund] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	fundInfo, err := dao.GetFundInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[PredictFund] [GetFundInfo] [err] = ", err.Error())
		return
	}

	fundData, err := dao.GetFundDataLimit30(c, fundInfo.TsCode)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[PredictFund] [GetFundData] [err] = ", err.Error())
		return
	}

	if len(fundData) == 0 {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[PredictFund] [GetFundData] [err] = 数据为空", "")
		return
	}

	sort.Slice(fundData, func(i, j int) bool {
		return fundData[i].TradeDate.Before(fundData[j].TradeDate)
	})

	last7 := make([]float64, 0, 7)
	for i := range fundData[:7] {
		last7 = append(last7, fundData[i].Close)
	}

	rdb := connector.GetRedis().WithContext(c)
	result, err := rdb.Get(c, fmt.Sprintf(public.RedisKeyFundPredict, req.Id)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[PredictFund] [rdb.Get] [err] = ", err.Error())
		return
	}
	if !errors.Is(err, redis.Nil) {
		util.SuccessResp(c, &PredictFundResp{
			List: last7,
			Val:  cast.ToFloat64(result),
		})
		return
	}

	val, err := python.PythonPredictFund(req.Id, fundData)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[PredictFund] [PythonPredict] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, &PredictFundResp{
		List: last7,
		Val:  val,
	})
}

func AiFund(c *gin.Context) {
	var req AiFundReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[AiFund] [ShouldBindJSON] [err] = %s", err.Error())
		return
	}

	fundInfo, err := dao.GetFundInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[AiFund] [GetFundInfo] [err] = %s", err.Error())
		return
	}

	limit30, err := dao.GetFundDataLimit30(c, fundInfo.TsCode)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[AiFund] [GetFundData] [err] = %s", err.Error())
		return
	}

	if len(limit30) == 0 {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[AiFund] [GetFundData] [err] = 数据为空", "")
		return
	}

	var close []float64
	for _, v := range limit30 {
		close = append(close, v.Close)
	}

	spark.SendSparkHttp(c, close, cast.ToString(util.GetUid(c)))
	return
}
