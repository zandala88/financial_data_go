package models

import (
	"context"
	"financia/public"
	"gorm.io/gorm"
	"time"
)

type StockForecast struct {
	Id      int64     `gorm:"primaryKey;autoIncrement;column:f_id;comment:'主键 ID，自增'" json:"id"`
	Company string    `gorm:"not null;column:f_company;comment:'公司名称'" json:"company"`
	Date    time.Time `gorm:"type:date;not null;column:f_date;comment:'日期'" json:"date"`
	Type    int       `gorm:"default:0;not null;column:f_type;comment:'预测类型，1：lstm模型预测值，2：arima模型预测值'" json:"type"`
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
