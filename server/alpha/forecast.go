package alpha

import (
	"context"
	"financia/models"
	"go.uber.org/zap"
)

var a = []float64{}

func InsertStockForecast(company string) {
	ctx := context.Background()
	// 1. 获取所有股票数据
	data := make([]*models.StockForecast, 0, len(a))
	stockRepo := models.NewStockRepo(ctx)
	offset30, err := stockRepo.DateOffset30(company)
	if err != nil {
		zap.S().Errorf("DateOffset30 Error: %v", err)
		return
	}
	zap.S().Debug(len(offset30))
	zap.S().Debug(len(a))
	for i := 0; i < len(a)-1; i++ {
		data = append(data, &models.StockForecast{
			Company: company,
			Date:    offset30[i].Date,
			Type:    1,
			Value:   a[i],
		})
	}
	data = append(data, &models.StockForecast{
		Company: company,
		Date:    getNextWeekday(),
		Type:    1,
		Value:   a[len(a)-1],
	})

	// 2. 插入数据
	forecastRepo := models.NewStockForecastRepo(ctx)
	err = forecastRepo.Insert(data)

}
