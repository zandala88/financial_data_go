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
	Real          map[string]float64 `json:"real"`
	Pred          map[string]float64 `json:"pred"`
	TomorrowClose float64            `json:"tomorrowClose"`
	TotalNum      int                `json:"totalNum"`
	TrueNum       int                `json:"trueNum"`
	Open          float64            `json:"open"`
	High          float64            `json:"high"`
	Low           float64            `json:"low"`
	Close         float64            `json:"close"`
	Volume        int64              `json:"volume"`
}

func GetStock(c *gin.Context) {
	// 获取 URL 参数
	var query GetStockReq
	err := c.ShouldBindQuery(&query)
	if err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetStock] 参数错误", err)
		return
	}

	// 查询真实值
	stockRepo := models.NewStockRepo(c)
	stockList, err := stockRepo.FindByCompanyAndDate(query.Name, query.Start, query.End)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetStock] stockRepo.FindByCompanyAndDate err = ", err)
		return
	}
	if len(stockList) == 0 {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetStock] len(stockList) == 0 | err = ", err)
		return
	}

	last := stockList[len(stockList)-1]
	resp := &GetStockResp{
		Real:   make(map[string]float64),
		Pred:   make(map[string]float64),
		Open:   last.Open,
		High:   last.High,
		Low:    last.Low,
		Close:  last.Close,
		Volume: last.Volume,
	}

	// 查询预测值
	stockForecastRepo := models.NewStockForecastRepo(c)
	stockForecastList, err := stockForecastRepo.FindByCompanyAndDate(query.Name, query.Start, query.End)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetStock] stockForecastRepo.FindByCompanyAndDate err = ", err)
		return
	}
	if len(stockForecastList) == 0 {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetStock] len(stockForecastList) == 0 | err = ", err)
		return
	}

	// slice to map
	for _, stock := range stockList {
		resp.Real[stock.Date.Format(time.DateOnly)] = stock.Close
	}
	for _, stockForecast := range stockForecastList {
		resp.Pred[stockForecast.Date.Format(time.DateOnly)] = stockForecast.Value
	}

	// todo 获取明日预测值
	tomorrow, err := stockForecastRepo.FindLastByCompany(query.Name)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetStock] stockForecastRepo.FindLastByCompany err = ", err)
		return
	}
	resp.TomorrowClose = tomorrow.Value + 2

	inPrice, err := stockForecastRepo.GetInPrice(query.Name, query.Start, query.End)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetStock] stockForecastRepo.GetInPrice err = ", err)
		return
	}
	resp.TotalNum = inPrice.TotalCount
	resp.TrueNum = inPrice.Times

	util.SuccessResp(c, resp)
}

type GetStockAllReq struct {
	Name  string `form:"name" binding:"required"`       // 必填字段
	Start string `form:"start" binding:"required,date"` // 必填字段，日期格式验证
	End   string `form:"end" binding:"required,date"`   // 必填字段，日期格式验证
}

type GetStockAllResp struct {
	Close map[string]float64 `json:"close"`
	Open  map[string]float64 `json:"open"`
	High  map[string]float64 `json:"high"`
	Low   map[string]float64 `json:"low"`
}

func GetStockAll(c *gin.Context) {
	// 获取 URL 参数
	var query GetStockAllReq
	err := c.ShouldBindQuery(&query)
	if err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetStockAll] 参数错误", err)
		return
	}

	stockRepo := models.NewStockRepo(c)
	stockList, err := stockRepo.FindByCompanyAndDate(query.Name, query.Start, query.End)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetStockAll] stockRepo.FindByCompanyAndDate err = ", err)
		return
	}
	if len(stockList) == 0 {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetStockAll] len(stockList) == 0 | err = ", err)
		return
	}

	resp := &GetStockAllResp{
		Close: make(map[string]float64, len(stockList)),
		Open:  make(map[string]float64, len(stockList)),
		High:  make(map[string]float64, len(stockList)),
		Low:   make(map[string]float64, len(stockList)),
	}

	for _, stock := range stockList {
		resp.Close[stock.Date.Format(time.DateOnly)] = stock.Close
		resp.Open[stock.Date.Format(time.DateOnly)] = stock.Open
		resp.High[stock.Date.Format(time.DateOnly)] = stock.High
		resp.Low[stock.Date.Format(time.DateOnly)] = stock.Low
	}

	util.SuccessResp(c, resp)
}

type GetCurrencyReq struct {
	Name  string `form:"name" binding:"required"`       // 必填字段
	Start string `form:"start" binding:"required,date"` // 必填字段，日期格式验证
	End   string `form:"end" binding:"required,date"`   // 必填字段，日期格式验证
}

type GetCurrencyResp struct {
	Real          map[string]float64 `json:"real"`
	Pred          map[string]float64 `json:"pred"`
	TomorrowClose float64            `json:"tomorrowClose"`
	TotalNum      int                `json:"totalNum"`
	TrueNum       int                `json:"trueNum"`
	Open          float64            `json:"open"`
	High          float64            `json:"high"`
	Low           float64            `json:"low"`
	Close         float64            `json:"close"`
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

	// 查询真实值
	currencyRepo := models.NewCurrencyRepo(c)
	currencyList, err := currencyRepo.FindByFromToAndDate(query.Name, query.Start, query.End)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetCurrency] currencyRepo.FindByFromToAndDate err = ", err)
		return
	}
	if len(currencyList) == 0 {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetCurrency] len(currencyList) == 0 | err = ", err)
		return
	}

	last := currencyList[len(currencyList)-1]
	resp := &GetCurrencyResp{
		Real:  make(map[string]float64, len(currencyList)),
		Pred:  make(map[string]float64, len(currencyList)),
		Open:  last.Open,
		High:  last.High,
		Low:   last.Low,
		Close: last.Close,
	}

	// 查询预测值
	currencyForecastRepo := models.NewCurrencyForecastRepo(c)
	currencyForecastList, err := currencyForecastRepo.FindByFromToAndDate(query.Name, query.Start, query.End)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetCurrency] currencyForecastRepo.FindByFromToAndDate err = ", err)
		return
	}
	if len(currencyForecastList) == 0 {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetCurrency] len(currencyForecastList) == 0 | err = ", err)
		return
	}

	// slice to map
	for _, currency := range currencyList {
		resp.Real[currency.Date.Format(time.DateOnly)] = currency.Close
	}
	for _, currencyForecast := range currencyForecastList {
		resp.Pred[currencyForecast.Date.Format(time.DateOnly)] = currencyForecast.Value
	}

	// todo 获取明日预测值
	tomorrow, err := currencyForecastRepo.FindLastBySymbol(query.Name)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetCurrency] currencyForecastRepo.FindLastByFromTo err = ", err)
		return
	}
	resp.TomorrowClose = tomorrow.Value

	inPrice, err := currencyForecastRepo.GetInPrice(query.Name, query.Start, query.End)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetCurrency] currencyForecastRepo.GetInPrice err = ", err)
		return
	}
	resp.TotalNum = inPrice.TotalCount
	resp.TrueNum = inPrice.Times

	util.SuccessResp(c, resp)
}

type GetCurrencyAllReq struct {
	Symbol string `form:"symbol" binding:"required"`     // 必填字段
	Start  string `form:"start" binding:"required,date"` // 必填字段，日期格式验证
	End    string `form:"end" binding:"required,date"`   // 必填字段，日期格式验证
}

type GetCurrencyAllResp struct {
	Close map[string]float64 `json:"close"`
	Open  map[string]float64 `json:"open"`
	High  map[string]float64 `json:"high"`
	Low   map[string]float64 `json:"low"`
}

func GetCurrencyAll(c *gin.Context) {
	// 获取 URL 参数
	var query GetCurrencyAllReq
	err := c.ShouldBindQuery(&query)
	if err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetCurrencyAll] 参数错误", err)
		return
	}

	currencyRepo := models.NewCurrencyRepo(c)
	currencyList, err := currencyRepo.FindByFromToAndDate(query.Symbol, query.Start, query.End)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetCurrencyAll] currencyRepo.FindByFromToAndDate err = ", err)
		return
	}
	if len(currencyList) == 0 {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetCurrencyAll] len(currencyList) == 0 | err = ", err)
		return
	}

	resp := &GetCurrencyAllResp{
		Close: make(map[string]float64, len(currencyList)),
		Open:  make(map[string]float64, len(currencyList)),
		High:  make(map[string]float64, len(currencyList)),
		Low:   make(map[string]float64, len(currencyList)),
	}

	for _, currency := range currencyList {
		resp.Close[currency.Date.Format(time.DateOnly)] = currency.Close
		resp.Open[currency.Date.Format(time.DateOnly)] = currency.Open
		resp.High[currency.Date.Format(time.DateOnly)] = currency.High
		resp.Low[currency.Date.Format(time.DateOnly)] = currency.Low
	}

	util.SuccessResp(c, resp)
}
