package fund

import (
	"financia/public/db/dao"
	"financia/server/tushare"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"time"
)

type DataFundReq struct {
	Id        int    `form:"id" binding:"required"`
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
}

type DataFundResp struct {
	Have bool              `json:"have"`
	List []*DataFundSimple `json:"list"`
}

type DataFundSimple struct {
	TradeDate string  `json:"tradeDate"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	PreClose  float64 `json:"preClose"`
	Change    float64 `json:"change"`
	PctChg    float64 `json:"pctChg"`
	Vol       int64   `json:"vol"`
	Amount    float64 `json:"amount"`
}

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
			TradeDate: v.TradeDate.Format("2006-01-02"),
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
		date := strings.ReplaceAll(last.TradeDate.Add(time.Hour*24).Format("2006-01-02"), "-", "")
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
