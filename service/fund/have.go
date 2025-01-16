package fund

import (
	"financia/public"
	"financia/public/db/dao"
	"financia/public/db/model"
	"financia/server/tushare"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HaveFundReq struct {
	Id int `form:"id" binding:"required"`
}

type HaveFundResp struct {
	Have bool `json:"have"`
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
		have = tushare.DailyFundAll(c, info.TsCode)
		if have {
			dao.UpdateFund(c, &model.FundInfo{
				Id:   info.Id,
				Flag: public.FundInfoFlagExist,
			})
		}
	}

	util.SuccessResp(c, &HaveFundResp{
		Have: have,
	})

}
