package service

import (
	"financia/models"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetInfoReq struct {
	Type int `form:"type" binding:"required"`
}

type GetInfoRespSimple struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type GetInfoResp struct {
	List []*GetInfoRespSimple `json:"list"`
}

func GetInfo(c *gin.Context) {
	var query GetInfoReq
	err := c.ShouldBindQuery(&query)
	if err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetInfo] 参数错误")
		return
	}

	alphaInfoRepo := models.NewAlphaInfoRepo(c)
	infoList, err := alphaInfoRepo.GetSymbolByType(query.Type)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	resp := &GetInfoResp{
		List: make([]*GetInfoRespSimple, 0, len(infoList)),
	}
	for _, info := range infoList {
		resp.List = append(resp.List, &GetInfoRespSimple{
			Name:   info.Name,
			Symbol: info.Symbol,
		})
	}

	util.SuccessResp(c, resp)
}

type GetDataIndexResp struct {
	Stock    GetDataIndexSimple `json:"stock"`
	Currency GetDataIndexSimple `json:"currency"`
}

type GetDataIndexSimple struct {
	MaxIncreaseValue GetDataIndexData `json:"maxIncreaseValue"`
	MaxIncreaseRate  GetDataIndexData `json:"maxIncreaseRate"`
	MaxDecreaseValue GetDataIndexData `json:"maxDecreaseValue"`
	MaxDecreaseRate  GetDataIndexData `json:"maxDecreaseRate"`
}

type GetDataIndexData struct {
	Value float64 `json:"value"`
	Name  string  `json:"name"`
}

func GetDataIndex(c *gin.Context) {
	// 获取股票数据
	stockRepo := models.NewStockRepo(c)
	stockInfo, err := stockRepo.FindDiffAndRatio()
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetDataIndex] stockRepo.FindDiffAndRatio err = ", err)
		return
	}

	stockIndex := GetDataIndexSimple{}
	for _, stock := range stockInfo {
		if stock.Diff > stockIndex.MaxIncreaseValue.Value {
			stockIndex.MaxIncreaseValue.Value = stock.Diff
			stockIndex.MaxIncreaseValue.Name = stock.Company
		}
		if stock.Ratio > stockIndex.MaxIncreaseRate.Value {
			stockIndex.MaxIncreaseRate.Value = stock.Ratio
			stockIndex.MaxIncreaseRate.Name = stock.Company
		}
		if stock.Diff < stockIndex.MaxDecreaseValue.Value {
			stockIndex.MaxDecreaseValue.Value = stock.Diff
			stockIndex.MaxDecreaseValue.Name = stock.Company
		}
		if stock.Ratio < stockIndex.MaxDecreaseRate.Value {
			stockIndex.MaxDecreaseRate.Value = stock.Ratio
			stockIndex.MaxDecreaseRate.Name = stock.Company
		}
	}

	// 获取外汇数据
	currencyRepo := models.NewCurrencyRepo(c)
	currencyInfo, err := currencyRepo.FindDiffAndRatio()
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[GetDataIndex] currencyRepo.FindDiffAndRatio err = ", err)
		return
	}

	currencyIndex := GetDataIndexSimple{}
	for _, currency := range currencyInfo {
		if currency.Diff > currencyIndex.MaxIncreaseValue.Value {
			currencyIndex.MaxIncreaseValue.Value = currency.Diff
			currencyIndex.MaxIncreaseValue.Name = currency.Symbol
		}
		if currency.Ratio > currencyIndex.MaxIncreaseRate.Value {
			currencyIndex.MaxIncreaseRate.Value = currency.Ratio
			currencyIndex.MaxIncreaseRate.Name = currency.Symbol
		}
		if currency.Diff < currencyIndex.MaxDecreaseValue.Value {
			currencyIndex.MaxDecreaseValue.Value = currency.Diff
			currencyIndex.MaxDecreaseValue.Name = currency.Symbol
		}
		if currency.Ratio < currencyIndex.MaxDecreaseRate.Value {
			currencyIndex.MaxDecreaseRate.Value = currency.Ratio
			currencyIndex.MaxDecreaseRate.Name = currency.Symbol
		}
	}

	resp := &GetDataIndexResp{
		Stock:    stockIndex,
		Currency: currencyIndex,
	}
	util.SuccessResp(c, resp)
}
