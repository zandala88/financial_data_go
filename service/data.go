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
	List   map[string]float64 `json:"list"`
	Open   float64            `json:"open"`
	High   float64            `json:"high"`
	Low    float64            `json:"low"`
	Close  float64            `json:"close"`
	Volume int64              `json:"volume"`
}

// todo 限制查看数据自选日期的范围，
func GetStock(c *gin.Context) {
	// 获取 URL 参数
	var query GetStockReq
	err := c.ShouldBindQuery(&query)
	if err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetStock] 参数错误", err)
		return
	}

	stockRepo := models.NewStockRepo(c)
	stockList, err := stockRepo.FindByCompanyAndDate(query.Name, query.Start, query.End)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetStock] 查询失败", err)
		return
	}
	// todo len(stockList) > 0 判断
	last := stockList[len(stockList)-1]
	resp := &GetStockResp{
		List:   make(map[string]float64),
		Open:   last.Open,
		High:   last.High,
		Low:    last.Low,
		Close:  last.Close,
		Volume: last.Volume,
	}
	for _, stock := range stockList {
		resp.List[stock.Date.Format(time.DateOnly)] = stock.Close
	}
	util.SuccessResp(c, resp)
}

func GetStockForecast(c *gin.Context) {

}

type GetCurrencyReq struct {
	From  string `form:"from" binding:"required"`       // 必填字段
	To    string `form:"to" binding:"required"`         // 必填字段
	Start string `form:"start" binding:"required,date"` // 必填字段，日期格式验证
	End   string `form:"end" binding:"required,date"`   // 必填字段，日期格式验证
}

type GetCurrencyResp struct {
	List  map[string]float64 `json:"list"`
	Open  float64            `json:"open"`
	High  float64            `json:"high"`
	Low   float64            `json:"low"`
	Close float64            `json:"close"`
}

func GetCurrency(c *gin.Context) {
	// 获取 URL 参数
	var query GetCurrencyReq
	err := c.ShouldBindQuery(&query)
	if err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetCurrency] 参数错误", err)
		return
	}

	currencyRepo := models.NewCurrencyRepo(c)
	currencyList, err := currencyRepo.FindByFromToAndDate(query.From, query.To, query.Start, query.End)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetCurrency] 查询失败", err)
		return
	}

	last := currencyList[len(currencyList)-1]
	resp := &GetCurrencyResp{
		List:  make(map[string]float64, len(currencyList)),
		Open:  last.Open,
		High:  last.High,
		Low:   last.Low,
		Close: last.Close,
	}
	for _, currency := range currencyList {
		resp.List[currency.Date.Format(time.DateOnly)] = currency.Close
	}

	util.SuccessResp(c, resp)
}
