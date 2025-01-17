package fund

import (
	"financia/server/tushare"
	"financia/util"
	"github.com/gin-gonic/gin"
	"sort"
)

type GraphFundResp struct {
	Radio []*tushare.FundSalesRatioResp `json:"radio"`
	Inst  []*tushare.FundSalesVolResp   `json:"inst"`
}

func GraphFund(c *gin.Context) {
	radio := tushare.FundSalesRatio(c)
	vol := tushare.FundSalesVol(c)

	sort.Slice(vol, func(i, j int) bool {
		// year sec
		if vol[i].Year == vol[j].Year {
			// quarter sec
			if vol[i].Quarter == vol[j].Quarter {
				return vol[i].Rank < vol[j].Rank
			}
		}
		return vol[i].Year < vol[j].Year
	})
	util.SuccessResp(c, &GraphFundResp{
		Radio: radio,
		Inst:  vol,
	})
}
