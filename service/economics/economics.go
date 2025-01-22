package economics

import (
	"financia/public"
	"financia/server/tushare"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sort"
)

func ShiborEconomics(c *gin.Context) {
	list := tushare.EconomicsShibor(c)
	sort.Slice(list, func(i, j int) bool {
		return list[i].Date < list[j].Date
	})
	util.SuccessResp(c, &ShiborEconomicsResp{
		List: list,
	})
}

func CnGdpEconomics(c *gin.Context) {
	var req CnGdpEconomicsReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[CnGdpEconomics] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	if req.Year > public.CnGdpEconomicsEndYear || req.Year < public.CnGdpEconomicsStartYear {
		util.FailRespWithCode(c, util.ReqDataError)
		zap.S().Error("[CnGdpEconomics] [ShouldBindJSON] [err] = ", "year is invalid")
		return
	}

	if req.Year == public.CnGdpEconomicsEndYear && req.Quarter == public.CnGdpEconomicsEndQuarter {
		util.FailRespWithCode(c, util.ReqDataError)
		zap.S().Error("[CnGdpEconomics] [ShouldBindJSON] [err] = ", "quarter is invalid")
		return
	}

	q := req.Year + "Q" + req.Quarter

	list := tushare.EconomicsCnGDP(c, q)

	util.SuccessResp(c, &CnGdpEconomicsResp{
		List: list,
	})
}

func CnCpiEconomics(c *gin.Context) {
	list := tushare.EconomicsCnCPI(c)
	sort.Slice(list, func(i, j int) bool {
		return list[i].Month < list[j].Month
	})
	util.SuccessResp(c, &CnCpiEconomicsResp{
		List: list,
	})
}
