package fund

import (
	"github.com/gin-gonic/gin"
)

type DataFundReq struct {
	Id        int    `form:"id" binding:"required"`
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
}

type DataFundResp struct {
	List []*DataFundSimple `json:"list"`
}

type DataFundSimple struct {
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

func DataFund(c *gin.Context) {
	// todo

}
