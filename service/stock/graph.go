package stock

import (
	"financia/public/db/dao"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GraphStockResp struct {
	IsHs     map[string]int `json:"isHs"`
	Exchange map[string]int `json:"exchange"`
	Market   map[string]int `json:"market"`
}

func GraphStock(c *gin.Context) {
	fields, err := dao.CountStockFields(c)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Errorf("DistinctStockFields error: %s", err.Error())
		return
	}

	util.SuccessResp(c, &GraphStockResp{
		IsHs:     fields["is_hs"],
		Exchange: fields["exchange"],
		Market:   fields["market"],
	})
}
