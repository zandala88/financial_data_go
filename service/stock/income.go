package stock

import (
	"financia/public/db/dao"
	"financia/server/tushare"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IncomeStockReq struct {
	Id int `form:"id" binding:"required"`
}

type IncomeStockResp struct {
	List []*tushare.StockIncomeResp `json:"list"`
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
