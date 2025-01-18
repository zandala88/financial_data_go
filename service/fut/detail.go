package fut

import (
	"financia/server/tushare"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sort"
)

// CU 铜
// SR 白糖
// CF 棉花
// AL 铝
// ZN 锌
// JD 鸡蛋
// FG 玻璃
// AP 苹果
// PP 聚丙烯
// RB 螺纹钢
type DetailFutReq struct {
	Prd string `form:"prd" binding:"required"`
}

type DetailFutResp struct {
	List []*tushare.FutWeeklyDetailResp `json:"list"`
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
