package models

import (
	"context"
	"financia/public"
	"gorm.io/gorm"
	"time"
)

type Crypto struct {
	Id       int64     `gorm:"primaryKey;autoIncrement;column:f_id;comment:'主键 ID，自增'" json:"id"`
	Currency string    `gorm:"not null;column:f_currency;comment:'货币类型'" json:"currency"`
	Date     time.Time `gorm:"type:date;not null;column:f_date;comment:'日期'" json:"date"`
	Open     float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_open;comment:'开盘价'" json:"open"`
	High     float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_high;comment:'最高价'" json:"high"`
	Low      float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_low;comment:'最低价'" json:"low"`
	Close    float64   `gorm:"type:decimal(10,3);default:0.000;not null;column:f_close;comment:'收盘价'" json:"close"`
	Volume   float64   `gorm:"default:0;not null;column:f_volume;comment:'成交量'" json:"volume"`
}

func (*Crypto) TableName() string {
	return "t_crypto_data"
}

type CryptoRepo struct {
	db  *gorm.DB
	ctx context.Context
}

func NewCryptoRepo(ctx context.Context) *CryptoRepo {
	return &CryptoRepo{
		db:  public.DB.WithContext(ctx),
		ctx: ctx,
	}
}

func (c *CryptoRepo) CreateCrypto(crypto *Crypto) error {
	err := c.db.Create(crypto).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *CryptoRepo) Insert(crypto []*Crypto) error {
	err := c.db.Create(crypto).Error
	if err != nil {
		return err
	}
	return nil
}
