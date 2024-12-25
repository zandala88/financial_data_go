package models

import (
	"context"
	"financia/public"
	"gorm.io/gorm"
	"time"
)

type Currency struct {
	Id    int64     `gorm:"primaryKey;autoIncrement;column:f_id;comment:'主键 ID，自增'" json:"id"`
	From  string    `gorm:"not null;column:f_from;comment:'货币来源'" json:"from"`
	To    string    `gorm:"not null;column:f_to;comment:'货币目标'" json:"to"`
	Date  time.Time `gorm:"type:date;not null;column:f_date;comment:'日期'" json:"date"`
	Open  float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_open;comment:'开盘价'" json:"open"`
	High  float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_high;comment:'最高价'" json:"high"`
	Low   float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_low;comment:'最低价'" json:"low"`
	Close float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_close;comment:'收盘价'" json:"close"`
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

func (c *CurrencyRepo) FindByFromToAndDate(from, to, start, end string) ([]*Currency, error) {
	// 根据货币来源、目标和日期范围查询
	var currency []*Currency
	err := c.db.Where("f_from = ? AND f_to = ? AND f_date >= ? AND f_date <= ?", from, to, start, end).Order("f_date").Find(&currency).Error
	if err != nil {
		return nil, err
	}
	return currency, nil
}
