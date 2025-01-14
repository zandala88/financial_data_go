package stock

import (
	"financia/public/db/dao"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type QueryStockResp struct {
	IsHsList     []string `json:"isHsList"`
	ExchangeList []string `json:"exchangeList"`
	MarketList   []string `json:"marketList"`
}

func QueryStock(c *gin.Context) {
	fields, err := dao.DistinctFields(c)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Errorf("DistinctFields error: %s", err.Error())
		return
	}

	util.SuccessResp(c, &QueryStockResp{
		IsHsList:     fields["is_hs"],
		ExchangeList: fields["exchange"],
		MarketList:   fields["market"],
	})
}
