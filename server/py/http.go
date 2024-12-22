package py

import (
	"context"
	"encoding/json"
	"financia/models"
	"financia/public"
	"financia/server/alpha"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
)

const (
	lstmUrl = "http://localhost:5000/predict"
)

type LSTM struct {
	Name string       `json:"name"`
	Data []*StockData `json:"data"`
}

type StockData struct {
	Open      float64 `json:"open"`
	Close     float64 `json:"close"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Volume    int64   `json:"volume"`
	NextClose float64 `json:"next_close"`
}

// todo 查询信息请求py接口
func Lstm() {
	ctx := context.Background()
	alphaInfoRepo := models.NewAlphaInfoRepo(ctx)
	stockList, err := alphaInfoRepo.GetSymbolByType(public.StockType)
	if err != nil {
		zap.S().Error("GetSymbolByType error: ", err)
		return
	}

	stockRepo := models.NewStockRepo(ctx)
	stockForecastRepo := models.NewStockForecastRepo(ctx)
	forecastList := make([]*models.StockForecast, 0, len(stockList))
	for _, stock := range stockList {
		stockData, err := stockRepo.FindLimitByCompany(30, stock.Symbol)
		if err != nil {
			zap.S().Error("FindLimitByCompany error: ", err)
			continue
		}

		var lstmData []*StockData
		for i := 0; i < len(stockData)-1; i++ {
			lstmData = append(lstmData, &StockData{
				Open:      stockData[i].Open,
				Close:     stockData[i].Close,
				High:      stockData[i].High,
				Low:       stockData[i].Low,
				Volume:    stockData[i].Volume,
				NextClose: stockData[i+1].Close,
			})
		}
		lstmData = append(lstmData, &StockData{
			Open:      stockData[len(stockData)-1].Open,
			Close:     stockData[len(stockData)-1].Close,
			High:      stockData[len(stockData)-1].High,
			Low:       stockData[len(stockData)-1].Low,
			Volume:    stockData[len(stockData)-1].Volume,
			NextClose: 0,
		})

		lstm := &LSTM{
			Name: stock.Symbol,
			Data: lstmData,
		}

		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(lstm).Post(lstmUrl)
		if err != nil {
			zap.S().Error("Post error: ", err)
			continue
		}
		zap.S().Info("Post Response: ", resp)
		respStruct := &alpha.LstmResp{}
		err = json.Unmarshal(resp.Body(), respStruct)
		if err != nil {
			zap.S().Error("Unmarshal error: ", err)
			continue
		}
		forecastList = append(forecastList, &models.StockForecast{
			Company: stock.Symbol,
			Date:    time.Now(),
			Type:    public.LstmType,
			Value:   respStruct.Predictions[0],
		})

	}
	err = stockForecastRepo.Insert(forecastList)
	if err != nil {
		zap.S().Error("Insert error: ", err)
		return
	}
	return
}

func LstmTest(name string) {
	ctx := context.Background()
	stockRepo := models.NewStockRepo(ctx)
	stockData, err := stockRepo.FindLimitByCompany(30, name)
	if err != nil {
		zap.S().Error("FindLimitByCompany error: ", err)
		return
	}

	var lstmData []*StockData
	for i := 0; i < len(stockData)-1; i++ {
		lstmData = append(lstmData, &StockData{
			Open:      stockData[i].Open,
			Close:     stockData[i].Close,
			High:      stockData[i].High,
			Low:       stockData[i].Low,
			Volume:    stockData[i].Volume,
			NextClose: stockData[i+1].Close,
		})
	}
	lstmData = append(lstmData, &StockData{
		Open:      stockData[len(stockData)-1].Open,
		Close:     stockData[len(stockData)-1].Close,
		High:      stockData[len(stockData)-1].High,
		Low:       stockData[len(stockData)-1].Low,
		Volume:    stockData[len(stockData)-1].Volume,
		NextClose: 0,
	})

	lstm := &LSTM{
		Name: "PDD",
		Data: lstmData,
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(lstm).Post(lstmUrl)
	if err != nil {
		zap.S().Error("Post error: ", err)
		return
	}
	zap.S().Info("Post Response: ", resp)
	respStruct := &alpha.LstmResp{}
	err = json.Unmarshal(resp.Body(), respStruct)
	if err != nil {
		zap.S().Error("Unmarshal error: ", err)
		return
	}
	fmt.Println(respStruct.Predictions[0])
}
