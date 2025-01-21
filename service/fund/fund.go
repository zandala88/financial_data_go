package fund

import (
	"financia/public/db/dao"
	"financia/server/tushare"
	"financia/util"
	"github.com/gin-gonic/gin"
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
			Vol:       v.Vol,
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

	util.SuccessResp(c, &DataFundResp{
		Have: true,
		List: respList,
	})
}

func GraphFund(c *gin.Context) {
	radio := tushare.FundSalesRatio(c)
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
	util.SuccessResp(c, &GraphFundResp{
		Radio: radio,
		Inst:  vol,
	})
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
