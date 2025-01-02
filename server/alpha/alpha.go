package alpha

import (
	"context"
	"encoding/json"
	"financia/config"
	"financia/models"
	"financia/public"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"slices"
	"time"
)

const (
	url = "https://www.alphavantage.co/query"

	dailyStockFuncName = "TIME_SERIES_DAILY"
	dailyCurrencyFunc  = "FX_DAILY"
	dailyCryptoFunc    = "DIGITAL_CURRENCY_DAILY"

	market  = "EUR"
	full    = "full"
	compact = "compact"
)

func AlphaDaily() {
	ctx := context.Background()

	// 查询type 1 股票参数名称
	alphaRepo := models.NewAlphaInfoRepo(ctx)
	alphaStockInfos, err := alphaRepo.GetSymbolByType(public.StockType)
	if err != nil {
		zap.S().Errorf("GetAlphaStock Error: %v", err)
	}
	stockList := make([]*models.Stock, 0, len(alphaStockInfos))
	for _, alphaInfo := range alphaStockInfos {
		alpha, err := GetAlphaStock(alphaInfo.Symbol)
		if err != nil {
			zap.S().Errorf("GetAlphaStock Error: %v", err)
			continue
		}
		stockList = append(stockList, alpha)
	}
	stockRepo := models.NewStockRepo(ctx)
	err = stockRepo.Insert(stockList)
	if err != nil {
		zap.S().Errorf("Insert Error: %v", err)
	}

	// 查询type 2 外汇参数名称
	alphaCurrencyInfos, err := alphaRepo.GetSymbolByType(public.CurrencyType)
	if err != nil {
		zap.S().Errorf("GetAlphaCurrency Error: %v", err)
	}
	currencyList := make([]*models.Currency, 0, len(alphaCurrencyInfos))
	for _, alphaInfo := range alphaCurrencyInfos {
		alpha, err := GetAlphaCurrency(alphaInfo.Symbol, alphaInfo.SymbolTo)
		if err != nil {
			zap.S().Errorf("GetAlphaCurrency Error: %v", err)
			continue
		}
		currencyList = append(currencyList, alpha)
	}
	currencyRepo := models.NewCurrencyRepo(ctx)
	err = currencyRepo.Insert(currencyList)
	if err != nil {
		zap.S().Errorf("Insert Error: %v", err)
	}

	// 查询type 3 加密货币参数名称
	alphaCryptoInfos, err := alphaRepo.GetSymbolByType(public.CryptoType)
	if err != nil {
		zap.S().Errorf("GetAlphaCrypto Error: %v", err)
	}
	cryptoList := make([]*models.Crypto, 0, len(alphaCryptoInfos))
	for _, alphaInfo := range alphaCryptoInfos {
		alpha, err := GetAlphaCrypto(alphaInfo.Symbol)
		if err != nil {
			zap.S().Errorf("GetAlphaCrypto Error: %v", err)
			continue
		}
		cryptoList = append(cryptoList, alpha)
	}
	cryptoRepo := models.NewCryptoRepo(ctx)
	err = cryptoRepo.Insert(cryptoList)
	if err != nil {
		zap.S().Errorf("Insert Error: %v", err)
	}

	return
}

func GetAlphaStock(symbol string) (*models.Stock, error) {
	client := resty.New()
	zap.S().Debugf("GetAlphaStock symbol: %s", symbol)
	get, err := client.R().SetQueryParam("function", dailyStockFuncName).
		SetQueryParam("symbol", symbol).
		SetQueryParam("outputsize", compact).
		SetQueryParam("apikey", config.Configs.Alpha.ApiKey).
		Get(url)
	if err != nil {
		zap.S().Errorf("GetAlphaStock Error: %v", err)
		return nil, err
	}
	resp := &StockResp{}
	err = json.Unmarshal(get.Body(), resp)
	if err != nil {
		zap.S().Errorf("GetAlphaStock Error: %v", err)
		return nil, err
	}

	yesterdayData, ok := resp.TimeSeriesDaily[getYesterdayStr()]
	zap.S().Debugf("GetAlphaStock yesterdayData: %+v", yesterdayData)
	if !ok {
		zap.S().Debug("No data for yesterday")
		return nil, nil
	}

	stock := &models.Stock{
		Company: resp.MetaData.Symbol,
		Date:    time.Now().AddDate(0, 0, -1),
		Open:    cast.ToFloat64(yesterdayData.Open),
		High:    cast.ToFloat64(yesterdayData.High),
		Low:     cast.ToFloat64(yesterdayData.Low),
		Close:   cast.ToFloat64(yesterdayData.Close),
		Volume:  cast.ToInt64(yesterdayData.Volume),
	}
	return stock, nil
}

func GetAlphaCurrency(fromSymbol, toSymbol string) (*models.Currency, error) {
	client := resty.New()
	get, err := client.R().SetQueryParam("function", dailyCurrencyFunc).
		SetQueryParam("from_symbol", fromSymbol).
		SetQueryParam("to_symbol", toSymbol).
		SetQueryParam("outputsize", compact).
		SetQueryParam("apikey", config.Configs.Alpha.ApiKey).
		Get(url)
	if err != nil {
		zap.S().Errorf("GetAlphaCurrency Error: %v", err)
		return nil, err
	}

	resp := &CurrencyResp{}
	err = json.Unmarshal(get.Body(), resp)
	if err != nil {
		zap.S().Errorf("GetAlphaCurrency Error: %v", err)
		return nil, err
	}

	yesterdayData, ok := resp.TimeSeriesDaily[getYesterdayStr()]
	zap.S().Debugf("GetAlphaCurrency yesterdayData: %+v", yesterdayData)
	if !ok {
		zap.S().Debug("No data for yesterday")
		return nil, nil
	}
	currency := &models.Currency{
		From:  fromSymbol,
		To:    toSymbol,
		Date:  time.Now().AddDate(0, 0, -1),
		Open:  cast.ToFloat64(yesterdayData.Open),
		High:  cast.ToFloat64(yesterdayData.High),
		Low:   cast.ToFloat64(yesterdayData.Low),
		Close: cast.ToFloat64(yesterdayData.Close),
	}
	return currency, nil
}

func GetAlphaCrypto(symbol string) (*models.Crypto, error) {
	client := resty.New()
	get, err := client.R().SetQueryParam("function", dailyCurrencyFunc).
		SetQueryParam("symbol", symbol).
		SetQueryParam("market", market).
		SetQueryParam("apikey", config.Configs.Alpha.ApiKey).
		Get(url)
	if err != nil {
		zap.S().Errorf("GetAlphaCurrency Error: %v", err)
		return nil, err
	}

	resp := &CryptoResp{}
	err = json.Unmarshal(get.Body(), resp)
	if err != nil {
		zap.S().Errorf("GetAlphaCurrency Error: %v", err)
		return nil, err
	}

	yesterdayData, ok := resp.TimeSeriesDaily[getYesterdayStr()]
	zap.S().Debugf("GetAlphaCurrency yesterdayData: %+v", yesterdayData)
	if !ok {
		zap.S().Debug("No data for yesterday")
		return nil, nil
	}

	crypto := &models.Crypto{
		Currency: symbol,
		Date:     time.Now().AddDate(0, 0, -1),
		Open:     cast.ToFloat64(yesterdayData.Open),
		High:     cast.ToFloat64(yesterdayData.High),
		Low:      cast.ToFloat64(yesterdayData.Low),
		Close:    cast.ToFloat64(yesterdayData.Close),
		Volume:   cast.ToFloat64(yesterdayData.Volume),
	}

	return crypto, nil
}

func InitInsertStockData(symbol string) {
	client := resty.New()
	get, err := client.R().SetQueryParam("function", dailyStockFuncName).
		SetQueryParam("symbol", symbol).
		SetQueryParam("outputsize", full).
		SetQueryParam("apikey", config.Configs.Alpha.ApiKey).
		Get(url)
	if err != nil {
		zap.S().Errorf("GetAlphaStock Error: %v", err)
		return
	}
	resp := &StockResp{}
	err = json.Unmarshal(get.Body(), resp)
	if err != nil {
		zap.S().Errorf("GetAlphaStock Error: %v", err)
		return
	}

	stockList := make([]*models.Stock, 0, len(resp.TimeSeriesDaily))
	for dateStr, data := range resp.TimeSeriesDaily {
		date, _ := time.Parse(time.DateOnly, dateStr)
		stock := &models.Stock{
			Company: resp.MetaData.Symbol,
			Date:    date,
			Open:    cast.ToFloat64(data.Open),
			High:    cast.ToFloat64(data.High),
			Low:     cast.ToFloat64(data.Low),
			Close:   cast.ToFloat64(data.Close),
			Volume:  cast.ToInt64(data.Volume),
		}
		stockList = append(stockList, stock)
	}
	slices.SortFunc(stockList, func(a, b *models.Stock) int {
		if a.Date.Before(b.Date) {
			return -1
		} else if a.Date.After(b.Date) {
			return 1 // a 在 b 之后
		}
		return 0
	})
	stockRepo := models.NewStockRepo(context.Background())
	err = stockRepo.Insert(stockList)
	if err != nil {
		zap.S().Errorf("Insert Error: %v", err)
		return
	}
	return
}

func InitInsertCurrencyData(fromSymbol, toSymbol string) {
	client := resty.New()
	get, err := client.R().SetQueryParam("function", dailyCurrencyFunc).
		SetQueryParam("from_symbol", fromSymbol).
		SetQueryParam("to_symbol", toSymbol).
		SetQueryParam("outputsize", full).
		SetQueryParam("apikey", config.Configs.Alpha.ApiKey).
		Get(url)
	if err != nil {
		zap.S().Errorf("GetAlphaStock Error: %v", err)
		return
	}
	resp := &CurrencyResp{}
	err = json.Unmarshal(get.Body(), resp)
	if err != nil {
		zap.S().Errorf("GetAlphaStock Error: %v", err)
		return
	}

	currencyList := make([]*models.Currency, 0, len(resp.TimeSeriesDaily))
	for dateStr, data := range resp.TimeSeriesDaily {
		date, _ := time.Parse(time.DateOnly, dateStr)
		stock := &models.Currency{
			From:  fromSymbol,
			To:    toSymbol,
			Date:  date,
			Open:  cast.ToFloat64(data.Open),
			High:  cast.ToFloat64(data.High),
			Low:   cast.ToFloat64(data.Low),
			Close: cast.ToFloat64(data.Close),
		}
		currencyList = append(currencyList, stock)
	}
	slices.SortFunc(currencyList, func(a, b *models.Currency) int {
		if a.Date.Before(b.Date) {
			return -1
		} else if a.Date.After(b.Date) {
			return 1 // a 在 b 之后
		}
		return 0
	})
	currencyRepo := models.NewCurrencyRepo(context.Background())
	err = currencyRepo.Insert(currencyList)
	if err != nil {
		zap.S().Errorf("Insert Error: %v", err)
		return
	}
	return
}

func InitInsertCryptoData(symbol string) {
	client := resty.New()
	get, err := client.R().SetQueryParam("function", dailyCryptoFunc).
		SetQueryParam("symbol", symbol).
		SetQueryParam("market", market).
		SetQueryParam("apikey", config.Configs.Alpha.ApiKey).
		Get(url)
	if err != nil {
		zap.S().Errorf("GetAlphaStock Error: %v", err)
		return
	}
	resp := &CryptoResp{}
	err = json.Unmarshal(get.Body(), resp)
	if err != nil {
		zap.S().Errorf("GetAlphaStock Error: %v", err)
		return
	}

	cryptoList := make([]*models.Crypto, 0, len(resp.TimeSeriesDaily))
	for dateStr, data := range resp.TimeSeriesDaily {
		date, _ := time.Parse(time.DateOnly, dateStr)
		stock := &models.Crypto{
			Currency: symbol,
			Date:     date,
			Open:     cast.ToFloat64(data.Open),
			High:     cast.ToFloat64(data.High),
			Low:      cast.ToFloat64(data.Low),
			Close:    cast.ToFloat64(data.Close),
			Volume:   cast.ToFloat64(data.Volume),
		}
		cryptoList = append(cryptoList, stock)
	}
	slices.SortFunc(cryptoList, func(a, b *models.Crypto) int {
		if a.Date.Before(b.Date) {
			return -1
		} else if a.Date.After(b.Date) {
			return 1 // a 在 b 之后
		}
		return 0
	})
	cryptoRepo := models.NewCryptoRepo(context.Background())
	err = cryptoRepo.Insert(cryptoList)
	if err != nil {
		zap.S().Errorf("Insert Error: %v", err)
		return
	}
	return
}
