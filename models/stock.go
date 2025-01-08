package models

import (
	"context"
	"financia/public"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Stock struct {
	Id      int64     `gorm:"primaryKey;autoIncrement;column:f_id;comment:'主键 ID，自增'" json:"id"`
	Company string    `gorm:"not null;column:f_company;comment:'公司名称'" json:"company"`
	Date    time.Time `gorm:"type:date;not null;column:f_date;comment:'日期'" json:"date"`
	Open    float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_open;comment:'开盘价'" json:"open"`
	High    float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_high;comment:'最高价'" json:"high"`
	Low     float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_low;comment:'最低价'" json:"low"`
	Close   float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_close;comment:'收盘价'" json:"close"`
	Volume  int64     `gorm:"default:0;not null;column:f_volume;comment:'成交量'" json:"volume"`
}

func (*Stock) TableName() string {
	return "t_stock_data"
}

type StockRepo struct {
	db  *gorm.DB
	ctx context.Context
}

func NewStockRepo(ctx context.Context) *StockRepo {
	return &StockRepo{
		db:  public.DB.WithContext(ctx),
		ctx: ctx,
	}
}

func (s *StockRepo) CreateStock(stock *Stock) error {
	err := s.db.Create(stock).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *StockRepo) Insert(stock []*Stock) error {
	err := s.db.Create(stock).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *StockRepo) FindLimitByCompany(limit int, company string) ([]*Stock, error) {
	// 根据date排序，取出limit条数据
	var stock []*Stock
	err := s.db.Where("f_company = ?", company).Order("f_date desc").Limit(limit).Find(&stock).Error
	if err != nil {
		zap.S().Error("FindLimitByCompany error: ", err)
		return nil, err
	}
	return stock, nil
}

func (s *StockRepo) FindByCompanyAndDate(company, start, end string) ([]*Stock, error) {
	// 根据公司名称和日期范围查询
	var stock []*Stock
	err := s.db.Where("f_company = ? AND f_date >= ? AND f_date <= ?", company, start, end).Order("f_date").Find(&stock).Error
	if err != nil {
		zap.S().Error("FindByCompanyAndDate error: ", err)
		return nil, err
	}
	return stock, nil
}

func (s *StockRepo) DateOffset30(company string) ([]*Stock, error) {
	var date []*Stock
	err := s.db.Model(&Stock{}).Where("f_company = ?", company).
		Select("f_date").Order("f_date").Limit(1000000).Offset(30).Find(&date).Error
	if err != nil {
		zap.S().Error("DateLimit30 error: ", err)
		return nil, err
	}
	return date, nil
}

type StockDiffAndRatioResult struct {
	Company string  `gorm:"column:f_company"`
	Diff    float64 `gorm:"column:close_diff"`
	Ratio   float64 `gorm:"column:close_ratio"`
}

func (s *StockRepo) FindDiffAndRatio() ([]*StockDiffAndRatioResult, error) {
	var result []*StockDiffAndRatioResult
	sql := `
WITH RankedStockData AS (
    SELECT
        f_company,
        f_date,
        f_close,
        ROW_NUMBER() OVER (PARTITION BY f_company ORDER BY f_date DESC) AS row_num
    FROM t_stock_data
),
     FilteredData AS (
         SELECT
             f_company,
             MAX(CASE WHEN row_num = 1 THEN f_close END) AS close_latest,
             MAX(CASE WHEN row_num = 2 THEN f_close END) AS close_previous
         FROM RankedStockData
         WHERE row_num <= 2
         GROUP BY f_company
     )
SELECT
    f_company,
    close_latest - close_previous AS close_diff,
    (close_latest - close_previous) / close_previous AS close_ratio
FROM FilteredData;`
	err := s.db.Raw(sql).Scan(&result).Error
	if err != nil {
		zap.S().Error("FindDiffAndRatio error: ", err)
		return nil, err
	}
	return result, nil
}
