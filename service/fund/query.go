package fund

import (
	"financia/public/db/dao"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type QueryFundResp struct {
	FundTypeList   []string `json:"fundTypeList"`
	InvestTypeList []string `json:"investTypeList"`
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
