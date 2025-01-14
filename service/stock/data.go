package stock

import (
	"financia/public/db/dao"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DataStockReq struct {
	Id        int    `form:"id" binding:"required"`
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
}

type DataStockResp struct {
	Have bool               `json:"have"`
	List []*DataStockSimple `json:"list"`
}

type DataStockSimple struct {
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
		util.SuccessResp(c, &DataStockResp{
			Have: false,
			List: nil,
		})
		return
	}

	respList := make([]*DataStockSimple, 0, len(list))
	for _, v := range list {
		respList = append(respList, &DataStockSimple{
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

	util.SuccessResp(c, &DataStockResp{
		Have: true,
		List: respList,
	})

}
