package fund

import (
	"encoding/json"
	"errors"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/dao"
	"financia/public/db/model"
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
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[DataFund] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	zap.S().Debugf("[DataFund] [req] = %#v", req)

	info, err := dao.GetFundInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[DataFund] [GetFundInfo] [err] = ", err.Error())
		return
	}

	list, err := dao.GetFundData(c, info.TsCode, req.StartDate, req.EndDate)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[DataFund] [GetFundData] [err] = ", err.Error())
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
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[DataFund] [rdb.SIsMember] [err] = ", err.Error())
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
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GraphFund] [rdb.Get] [err] = ", err.Error())
		return
	}
	if errors.Is(err, redis.Nil) {
		radio := tushare.FundSalesRatio(c)

		go func() {
			listStr, _ := json.Marshal(radio)
			_, err := rdb.Set(c, public.RedisKeyFundRadio, listStr, time.Duration(util.SecondsUntilMidnight())*time.Second).Result()
			if err != nil {
				util.FailRespWithCode(c, util.InternalServerError)
				zap.S().Error("[GraphFund] [rdb.Set] [err] = ", err.Error())
				return
			}
		}()

		resp.Radio = radio
	} else {
		var radio []*tushare.FundSalesRatioResp
		if err := json.Unmarshal([]byte(radioResult), &radio); err != nil {
			util.FailRespWithCode(c, util.InternalServerError)
			zap.S().Error("[GraphFund] [json.Unmarshal] [err] = ", err.Error())
			return
		}
		resp.Radio = radio
	}

	volResult, err := rdb.Get(c, public.RedisKeyFundVol).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GraphFund] [rdb.Get] [err] = ", err.Error())
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
			_, err := rdb.Set(c, public.RedisKeyFundVol, listStr, time.Duration(util.SecondsUntilMidnight())*time.Second).Result()
			if err != nil {
				util.FailRespWithCode(c, util.InternalServerError)
				zap.S().Error("[GraphFund] [rdb.Set] [err] = ", err.Error())
				return
			}
		}()

		resp.Inst = vol
	} else {
		var vol []*tushare.FundSalesVolResp
		if err := json.Unmarshal([]byte(volResult), &vol); err != nil {
			util.FailRespWithCode(c, util.InternalServerError)
			zap.S().Error("[GraphFund] [json.Unmarshal] [err] = ", err.Error())
			return
		}
		resp.Inst = vol
	}

	util.SuccessResp(c, resp)
}

func HaveFund(c *gin.Context) {
	var req HaveFundReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[DataFund] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	info, err := dao.GetFundInfo(c, req.Id)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[DataFund] [GetFundInfo] [err] = ", err.Error())
		return
	}

	have, err := dao.CheckFundData(c, info.TsCode)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[DataFund] [GetFundData] [err] = ", err.Error())
		return
	}

	if !have {
		data := tushare.DailyFundAll(c, &tushare.DailyReq{
			TsCode: info.TsCode,
		})
		have = len(data) > 0
		if err := dao.InsertFundData(c, data); err != nil {
			util.FailRespWithCode(c, util.InternalServerError)
			zap.S().Errorf("[Daily] [InsertFundData] [err] = %s", err.Error())
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
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[ListStock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	list, count, err := dao.GetFundList(c, req.Search, req.FundType, req.InvestType, req.Page, req.PageSize)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[ListStock] [GetStockList] [err] = ", err.Error())
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
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Errorf("DistinctFundFields error: %s", err.Error())
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
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[FollowFund] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	userId := util.GetUid(c)
	rdb := connector.GetRedis().WithContext(c)
	redisKey := fmt.Sprintf(public.RedisKeyFundFollow, userId)

	exists, err := rdb.SIsMember(c, redisKey, req.Id).Result()
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[FollowFund] [rdb.SIsMember] [err] = ", err.Error())
		return
	}
	if req.Follow == exists {
		util.SuccessResp(c, nil)
		return
	}

	if req.Follow {
		_, err = rdb.SAdd(c, redisKey, req.Id).Result()
		if err != nil {
			util.FailRespWithCode(c, util.InternalServerError)
			zap.S().Error("[FollowFund] [rdb.SAdd] [err] = ", err.Error())
			return
		}
	} else {
		_, err = rdb.SRem(c, redisKey, req.Id).Result()
		if err != nil {
			util.FailRespWithCode(c, util.InternalServerError)
			zap.S().Error("[FollowFund] [rdb.SRem] [err] = ", err.Error())
			return
		}
	}

	util.SuccessResp(c, nil)
}
