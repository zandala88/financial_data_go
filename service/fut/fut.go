package fut

import (
	"encoding/json"
	"financia/public/db/connector"
	"financia/server/tushare"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sort"
	"time"
)

func QueryFut(c *gin.Context) {
	util.SuccessResp(c, &QueryFutResp{
		List: []*QueryFutSimple{
			{Prd: "CU", Name: "铜"},
			{Prd: "SR", Name: "白糖"},
			{Prd: "CF", Name: "棉花"},
			{Prd: "AL", Name: "铝"},
			{Prd: "ZN", Name: "锌"},
			{Prd: "JD", Name: "鸡蛋"},
			{Prd: "FG", Name: "玻璃"},
			{Prd: "AP", Name: "苹果"},
			{Prd: "PP", Name: "聚丙烯"},
			{Prd: "RB", Name: "螺纹钢"},
			{Prd: "RO", Name: "菜籽油"},
			{Prd: "M", Name: "豆粕"},
			{Prd: "JM", Name: "焦煤"},
			{Prd: "ZC", Name: "动力煤"},
			{Prd: "Y", Name: "豆油"},
			{Prd: "SS", Name: "不锈钢"},
			{Prd: "BU", Name: "沥青"},
			{Prd: "C", Name: "玉米"},
			{Prd: "AU", Name: "黄金"},
			{Prd: "RU", Name: "天胶"},
			{Prd: "RR", Name: "粳米"},
			{Prd: "RS", Name: "油菜籽"},
		},
	})
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

func DetailFut(c *gin.Context) {
	var req DetailFutReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[DetailFut] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	list := tushare.FutWeeklyDetail(c, req.Prd)
	sort.Slice(list, func(i, j int) bool {
		return list[i].WeekDate < list[j].WeekDate
	})

	util.SuccessResp(c, &DetailFutResp{
		List: list,
	})
}
