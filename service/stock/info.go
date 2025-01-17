package stock

import (
	"financia/public/db/dao"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InfoStockReq struct {
	Id int `form:"id" binding:"required"`
}

type InfoStockResp struct {
	FullName string `json:"name"`
	Industry string `json:"industry"`
	Market   string `json:"market"`
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

	util.SuccessResp(c, &InfoStockResp{
		FullName: info.FullName,
		Industry: info.Industry,
		Market:   info.Market,
	})
}
