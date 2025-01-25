package server

import (
	"encoding/csv"
	"encoding/json"
	"financia/config"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"os"
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

type StockResp struct {
	MetaData        StockMetaData        `json:"Meta Data"`
	TimeSeriesDaily StockTimeSeriesDaily `json:"Time Series (Daily)"`
}

type StockTimeSeriesDaily map[string]StockDailyData

type StockMetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

type StockDailyData struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type TmpStockData struct {
	Date   time.Time `json:"date"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume int64     `json:"volume"`
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

	stockList := make([]*TmpStockData, 0, len(resp.TimeSeriesDaily))
	for dateStr, data := range resp.TimeSeriesDaily {
		date, _ := time.Parse(time.DateOnly, dateStr)
		stock := &TmpStockData{
			Date:   date,
			Open:   cast.ToFloat64(data.Open),
			High:   cast.ToFloat64(data.High),
			Low:    cast.ToFloat64(data.Low),
			Close:  cast.ToFloat64(data.Close),
			Volume: cast.ToInt64(data.Volume),
		}
		stockList = append(stockList, stock)
	}
	slices.SortFunc(stockList, func(a, b *TmpStockData) int {
		if a.Date.Before(b.Date) {
			return -1
		} else if a.Date.After(b.Date) {
			return 1 // a 在 b 之后
		}
		return 0
	})
	TT1(stockList)
}

func TT1(stockList []*TmpStockData) {
	file, err := os.Create("tmp.csv")
	if err != nil {
		zap.S().Errorf("Failed to create CSV file: %v", err)
		return
	}
	defer file.Close()

	// 创建 CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Date", "Open", "High", "Low", "Volume", "Close"})
	if err != nil {
		zap.S().Errorf("Failed to write CSV header: %v", err)
		return
	}

	for _, stock := range stockList {
		row := []string{
			stock.Date.Format("2006-01-02"), // 日期格式化
			fmt.Sprintf("%.2f", stock.Open),
			fmt.Sprintf("%.2f", stock.High),
			fmt.Sprintf("%.2f", stock.Low),
			fmt.Sprintf("%d", stock.Volume),
			fmt.Sprintf("%.2f", stock.Close),
		}
		err := writer.Write(row)
		if err != nil {
			zap.S().Errorf("Failed to write CSV row: %v", err)
			return
		}
	}

}
