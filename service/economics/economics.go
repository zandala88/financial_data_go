package economics

import (
	"encoding/json"
	"errors"
	"financia/public"
	"financia/public/db/connector"
	"financia/server/tushare"
	"financia/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"sort"
)

func ShiborEconomics(c *gin.Context) {
	rdb := connector.GetRedis().WithContext(c)
	result, err := rdb.Get(c, public.RedisKeyShiborEconomics).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[ShiborEconomics] [rdb.Get] [err] = ", err.Error())
		return
	}
	if errors.Is(err, redis.Nil) {
		list := tushare.EconomicsShibor(c)
		sort.Slice(list, func(i, j int) bool {
			return list[i].Date < list[j].Date
		})

		go func() {
			listStr, _ := json.Marshal(list)
			if _, err := rdb.Set(c, public.RedisKeyShiborEconomics, listStr, 0).Result(); err != nil {
				util.FailRespWithCodeAndZap(c, util.InternalServerError, "[ShiborEconomics] [rdb.Set] [err] = ", err.Error())
				return
			}
		}()

		util.SuccessResp(c, &ShiborEconomicsResp{
			List: list,
		})
		return
	}

	var list []*tushare.EconomicsShiborResp
	if err := json.Unmarshal([]byte(result), &list); err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[ShiborEconomics] [json.Unmarshal] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, &ShiborEconomicsResp{
		List: list,
	})
}

func CnGdpEconomics(c *gin.Context) {
	var req CnGdpEconomicsReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[CnGdpEconomics] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	if req.Year > public.CnGdpEconomicsEndYear || req.Year < public.CnGdpEconomicsStartYear {
		util.FailRespWithCodeAndZap(c, util.ReqDataError, "[CnGdpEconomics] [ShouldBindJSON] [err] = ", "year is invalid")
		return
	}

	if req.Year == public.CnGdpEconomicsEndYear && req.Quarter == public.CnGdpEconomicsEndQuarter {
		util.FailRespWithCodeAndZap(c, util.ReqDataError, "[CnGdpEconomics] [ShouldBindJSON] [err] = ", "quarter is invalid")
		return
	}

	q := req.Year + "Q" + req.Quarter

	rdb := connector.GetRedis().WithContext(c)
	result, err := rdb.Get(c, public.RedisKeyCnGdpEconomics+q).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[CnGdpEconomics] [rdb.Get] [err] = ", err.Error())
		return
	}
	if errors.Is(err, redis.Nil) {
		list := tushare.EconomicsCnGDP(c, q)

		go func() {
			listStr, _ := json.Marshal(list)
			if _, err := rdb.Set(c, public.RedisKeyCnGdpEconomics+q, listStr, 0).Result(); err != nil {
				zap.S().Error("[CnGdpEconomics] [rdb.Set] [err] = ", err.Error())
				return
			}
		}()

		util.SuccessResp(c, &CnGdpEconomicsResp{
			List: list,
		})
		return
	}

	var list []*tushare.EconomicsCnGDPResp
	if err := json.Unmarshal([]byte(result), &list); err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[CnGdpEconomics] [json.Unmarshal] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, &CnGdpEconomicsResp{
		List: list,
	})
}

func CnCpiEconomics(c *gin.Context) {
	rdb := connector.GetRedis().WithContext(c)
	result, err := rdb.Get(c, public.RedisKeyCnCpiEconomics).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[CnCpiEconomics] [rdb.Get] [err] = ", err.Error())
		return
	}
	if errors.Is(err, redis.Nil) {
		list := tushare.EconomicsCnCPI(c)
		sort.Slice(list, func(i, j int) bool {
			return list[i].Month < list[j].Month
		})

		go func() {
			listStr, _ := json.Marshal(list)
			_, err := rdb.Set(c, public.RedisKeyCnCpiEconomics, listStr, 0).Result()
			if err != nil {
				zap.S().Error("[CnCpiEconomics] [rdb.Set] [err] = ", err.Error())
				return
			}
		}()

		util.SuccessResp(c, &CnCpiEconomicsResp{
			List: list,
		})
		return
	}

	var list []*tushare.EconomicsCnCPIResp
	if err := json.Unmarshal([]byte(result), &list); err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[CnCpiEconomics] [json.Unmarshal] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, &CnCpiEconomicsResp{
		List: list,
	})
}
