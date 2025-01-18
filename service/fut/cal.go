package fut

import (
	"encoding/json"
	"financia/public/db/connector"
	"financia/server/tushare"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type CalFutResp struct {
	Sse  []*tushare.FutTradeCalResp `json:"sse"`
	Szse []*tushare.FutTradeCalResp `json:"szse"`
}

func CalFut(c *gin.Context) {
	resp := &CalFutResp{
		Sse:  make([]*tushare.FutTradeCalResp, 0),
		Szse: make([]*tushare.FutTradeCalResp, 0),
	}

	rdb := connector.GetRedis().WithContext(c)
	result, err := rdb.Get(c, "cal_fut").Result()
	if err != nil {
		zap.S().Errorf("[CalFut] [rdb.Get] [err] = %s", err.Error())
		tuShareSet(c, resp)
		// 当天0点过期
		exp := util.SecondsUntilMidnight()
		rdbStr, _ := json.Marshal(resp)
		_, err := rdb.Set(c, "cal_fut", rdbStr, time.Duration(exp)*time.Second).Result()
		if err != nil {
			zap.S().Errorf("[CalFut] [rdb.Set] [err] = %s", err.Error())
		}
	} else {
		err = json.Unmarshal([]byte(result), resp)
		if err != nil {
			util.FailRespWithCode(c, util.InternalServerError)
			zap.S().Errorf("[CalFut] [json.Unmarshal] [err] = %s", err.Error())
			tuShareSet(c, resp)
		}
	}

	util.SuccessResp(c, resp)
}

func tuShareSet(c *gin.Context, resp *CalFutResp) {
	resp.Sse, resp.Szse = tushare.FutTradeCal(c)

	for i, v := range resp.Sse {
		resp.Sse[i].CalDate = util.ConvertDateStrToTime(v.CalDate, "20060102").Format(time.DateOnly)
	}
	for i, v := range resp.Szse {
		resp.Szse[i].CalDate = util.ConvertDateStrToTime(v.CalDate, "20060102").Format(time.DateOnly)
	}
}
