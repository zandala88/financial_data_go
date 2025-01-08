package models

import (
	"context"
	"financia/public"
	"gorm.io/gorm"
	"time"
)

type CurrencyForecast struct {
	Id     int64     `gorm:"primaryKey;autoIncrement;column:f_id;comment:'主键 ID，自增'" json:"id"`
	Symbol string    `gorm:"default:'';not null;column:f_symbol;comment:'货币符号'" json:"symbol"`
	Date   time.Time `gorm:"type:date;not null;column:f_date;comment:'日期'" json:"date"`
	Value  float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_value;comment:'预测值'" json:"value"`
}

func (*CurrencyForecast) TableName() string {
	return "t_currency_forecast_data"
}

type CurrencyForecastRepo struct {
	db  *gorm.DB
	ctx context.Context
}

func NewCurrencyForecastRepo(ctx context.Context) *CurrencyForecastRepo {
	return &CurrencyForecastRepo{
		db:  public.DB.WithContext(ctx),
		ctx: ctx,
	}
}

func (c *CurrencyForecastRepo) CreateCurrencyForecast(currencyForecast *CurrencyForecast) error {
	err := c.db.Create(currencyForecast).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CurrencyForecastRepo) Insert(currencyForecast []*CurrencyForecast) error {
	err := c.db.Create(currencyForecast).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CurrencyForecastRepo) FindByFromToAndDate(symbol, start, end string) ([]*CurrencyForecast, error) {
	// 根据货币来源、目标和日期范围查询
	var currencyForecast []*CurrencyForecast
	err := c.db.Where("f_symbol = ? AND f_date >= ? AND f_date <= ?", symbol, start, end).Order("f_date").Find(&currencyForecast).Error
	if err != nil {
		return nil, err
	}
	return currencyForecast, nil
}

func (c *CurrencyForecastRepo) FindLastBySymbol(symbol string) (*CurrencyForecast, error) {
	var currencyForecast CurrencyForecast
	err := c.db.Where("f_symbol = ?", symbol).Order("f_date desc").First(&currencyForecast).Error
	if err != nil {
		return nil, err
	}
	return &currencyForecast, nil
}

func (c *CurrencyForecastRepo) GetInPrice(symbol, start, end string) (*GetInPriceResult, error) {
	var result GetInPriceResult
	sql := `
WITH ranked_data AS (
    SELECT
        f_close,
        f_value,
        t_currency_data.f_date,
        LEAD(f_close) OVER (PARTITION BY t_currency_data.f_symbol ORDER BY f_date) AS next_f_close,
        LEAD(f_value) OVER (PARTITION BY t_currency_data.f_symbol ORDER BY f_date) AS next_f_value,
        t_currency_data.f_symbol
    FROM t_currency_data
             JOIN t_currency_forecast_data
                  ON t_currency_data.f_symbol = t_currency_forecast_data.f_symbol
                      AND t_currency_data.f_date = t_currency_forecast_data.f_date
    WHERE t_currency_data.f_date BETWEEN ? AND ?
      AND t_currency_data.f_symbol = ?
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
FROM ranked_data;
`
	err := c.db.Raw(sql, start, end, symbol).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.Times < 0 {
		result.Times = 0
	}
	return &result, nil
}
