package fut

import (
	"financia/server/tushare"
	"financia/util"
	"github.com/gin-gonic/gin"
	"time"
)

func tuShareSet(c *gin.Context, resp *CalFutResp) {
	resp.Sse, resp.Szse = tushare.FutTradeCal(c)

	for i, v := range resp.Sse {
		resp.Sse[i].CalDate = util.ConvertDateStrToTime(v.CalDate, "20060102").Format(time.DateOnly)
	}
	for i, v := range resp.Szse {
		resp.Szse[i].CalDate = util.ConvertDateStrToTime(v.CalDate, "20060102").Format(time.DateOnly)
	}
}
