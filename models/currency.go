package models

import (
	"context"
	"financia/public"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Currency struct {
	Id     int64     `gorm:"primaryKey;autoIncrement;column:f_id;comment:'主键 ID，自增'" json:"id"`
	Symbol string    `gorm:"default:'';not null;column:f_symbol;comment:'货币符号'" json:"symbol"`
	Date   time.Time `gorm:"type:date;not null;column:f_date;comment:'日期'" json:"date"`
	Open   float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_open;comment:'开盘价'" json:"open"`
	High   float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_high;comment:'最高价'" json:"high"`
	Low    float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_low;comment:'最低价'" json:"low"`
	Close  float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_close;comment:'收盘价'" json:"close"`
}

func (*Currency) TableName() string {
	return "t_currency_data"
}

type CurrencyRepo struct {
	db  *gorm.DB
	ctx context.Context
}

func NewCurrencyRepo(ctx context.Context) *CurrencyRepo {
	return &CurrencyRepo{
		db:  public.DB.WithContext(ctx),
		ctx: ctx,
	}
}

func (c *CurrencyRepo) CreateCurrency(currency *Currency) error {
	err := c.db.Create(currency).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CurrencyRepo) Insert(currency []*Currency) error {
	err := c.db.Create(currency).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CurrencyRepo) FindByFromToAndDate(symbol, start, end string) ([]*Currency, error) {
	// 根据货币来源、目标和日期范围查询
	var currency []*Currency
	err := c.db.Where("f_symbol = ? AND f_date >= ? AND f_date <= ?", symbol, start, end).Order("f_date").Find(&currency).Error
	if err != nil {
		return nil, err
	}
	return currency, nil
}

type CurrencyDiffAndRatioResult struct {
	Symbol string  `gorm:"column:f_name"`
	Diff   float64 `gorm:"column:close_diff"`
	Ratio  float64 `gorm:"column:close_ratio"`
}

func (c *CurrencyRepo) FindDiffAndRatio() ([]*CurrencyDiffAndRatioResult, error) {
	var result []*CurrencyDiffAndRatioResult
	sql := `
WITH RankedStockData AS (
    SELECT
        f_symbol,
        f_date,
        f_close,
        ROW_NUMBER() OVER (PARTITION BY f_symbol ORDER BY f_date DESC) AS row_num
    FROM t_currency_data
),
FilteredData AS (
    SELECT
        f_symbol,
        MAX(CASE WHEN row_num = 1 THEN f_close END) AS close_latest,
        MAX(CASE WHEN row_num = 2 THEN f_close END) AS close_previous
    FROM RankedStockData
    WHERE row_num <= 2
    GROUP BY f_symbol
),
CalculatedData AS (
    SELECT
        f_symbol, close_latest - close_previous AS close_diff, (close_latest - close_previous) / close_previous AS close_ratio
    FROM FilteredData
)
SELECT
    ai.f_name,
    cd.close_diff,
    cd.close_ratio
FROM CalculatedData cd
         JOIN t_alpha_info ai
              ON cd.f_symbol = ai.f_symbol;
`
	err := c.db.Raw(sql).Scan(&result).Error
	if err != nil {
		zap.S().Error("FindDiffAndRatio error: ", err)
		return nil, err
	}
	return result, nil
}
