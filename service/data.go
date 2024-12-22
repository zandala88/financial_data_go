package service

import (
	"financia/models"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type GetStockReq struct {
	Name  string `form:"name" binding:"required"`       // 必填字段
	Start string `form:"start" binding:"required,date"` // 必填字段，日期格式验证
	End   string `form:"end" binding:"required,date"`   // 必填字段，日期格式验证
}

type GetStockResp struct {
	List map[string]float64 `json:"list"`
}

func GetStock(c *gin.Context) {
	// 获取 URL 参数
	var query GetStockReq
	err := c.ShouldBindQuery(&query)
	if err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetStock] 参数错误")
		return
	}

	stockRepo := models.NewStockRepo(c)
	stockList, err := stockRepo.FindByCompanyAndDate(query.Name, query.Start, query.End)
	if err != nil {
		util.FailResp(c, err.Error())
		return
	}
	resp := &GetStockResp{
		List: make(map[string]float64),
	}
	for _, stock := range stockList {
		resp.List[stock.Date.Format(time.DateOnly)] = stock.Close
	}
	util.SuccessResp(c, resp)
}

func GetStockForecast(c *gin.Context) {

}
