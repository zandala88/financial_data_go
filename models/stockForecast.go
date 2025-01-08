package models

import (
	"context"
	"financia/public"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type StockForecast struct {
	Id      int64     `gorm:"primaryKey;autoIncrement;column:f_id;comment:'主键 ID，自增'" json:"id"`
	Company string    `gorm:"not null;column:f_company;comment:'公司名称'" json:"company"`
	Date    time.Time `gorm:"type:date;not null;column:f_date;comment:'日期'" json:"date"`
	Value   float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_value;comment:'预测值'" json:"value"`
}

func (*StockForecast) TableName() string {
	return "t_stock_forecast_data"
}

type StockForecastRepo struct {
	db  *gorm.DB
	ctx context.Context
}

func NewStockForecastRepo(ctx context.Context) *StockForecastRepo {
	return &StockForecastRepo{
		db:  public.DB.WithContext(ctx),
		ctx: ctx,
	}
}

func (s *StockForecastRepo) CreateStockForecast(stock *StockForecast) error {
	err := s.db.Create(stock).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *StockForecastRepo) Insert(stock []*StockForecast) error {
	err := s.db.Create(stock).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *StockForecastRepo) FindByCompanyAndDate(company, start, end string) ([]*StockForecast, error) {
	// 根据公司名称和日期范围查询
	var stock []*StockForecast
	err := s.db.Where("f_company = ? AND f_date BETWEEN ? AND ?", company, start, end).Order("f_date").Find(&stock).Error
	if err != nil {
		zap.S().Error("FindByCompanyAndDate error: ", err)
		return nil, err
	}
	return stock, nil
}

func (s *StockForecastRepo) FindLastByCompany(company string) (*StockForecast, error) {
	// 根据公司名称查询最后一条数据
	var stock StockForecast
	err := s.db.Where("f_company = ?", company).Order("f_date desc").First(&stock).Error
	if err != nil {
		zap.S().Error("FindLastByCompany error: ", err)
		return nil, err
	}
	return &stock, nil
}

type GetInPriceResult struct {
	TotalCount int `gorm:"column:total_count"`
	Times      int `gorm:"column:times"`
}

func (s *StockForecastRepo) GetInPrice(company, start, end string) (*GetInPriceResult, error) {
	result := &GetInPriceResult{}
	sql := `WITH ranked_data AS (
    SELECT
        f_close,
        f_value,
        t_stock_data.f_date,
        LEAD(f_close) OVER (PARTITION BY t_stock_data.f_company ORDER BY f_date) AS next_f_close,
        LEAD(f_value) OVER (PARTITION BY t_stock_data.f_company ORDER BY f_date) AS next_f_value,
        t_stock_data.f_company
    FROM t_stock_data
             JOIN t_stock_forecast_data
                  ON t_stock_data.f_company = t_stock_forecast_data.f_company
                      AND t_stock_data.f_date = t_stock_forecast_data.f_date
    WHERE t_stock_data.f_date BETWEEN ? AND ?
      AND t_stock_data.f_company = ?
)
SELECT
    COUNT(1) AS total_count,
    SUM(
            CASE
                WHEN (next_f_close - f_close > 0 AND next_f_value - f_value > 0)
                    OR (next_f_close - f_close < 0 AND next_f_value - f_value < 0)
					OR (next_f_close - f_close = 0 AND next_f_value - f_value = 0)
                    THEN 1
                ELSE 0
                END
    ) - 1 AS times
FROM ranked_data;`

	err := s.db.Raw(sql, start, end, company).Scan(&result).Error
	if err != nil {
		zap.S().Error("GetInPrice error: ", err)
		return nil, err
	}
	if result.Times < 0 {
		result.Times = 0
	}
	return result, nil
}
