package stock

import (
	"financia/public/db/dao"
	"financia/server/tushare"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HaveStockReq struct {
	Id int `form:"id" binding:"required"`
}

type HaveStockResp struct {
	Have bool `json:"have"`
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
		tushare.DailyAll(c, info.TsCode)
		have = true
	}

	util.SuccessResp(c, &HaveStockResp{
		Have: true,
	})
	return
}
