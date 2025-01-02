package models

import (
	"context"
	"financia/public"
	"gorm.io/gorm"
	"time"
)

type AlphaInfo struct {
	Id        int64      `gorm:"primaryKey;autoIncrement;column:f_id;comment:'主键 ID'" json:"id"`
	Name      string     `gorm:"default:'';not null;column:f_name;comment:'名称'" json:"name"`
	Symbol    string     `gorm:"default:'';not null;column:f_symbol;comment:'Alpha Vantage API 参数'" json:"symbol"`
	SymbolTo  string     `gorm:"default:'';not null;column:f_symbol_to;comment:'外汇参数'" json:"symbol_to"`
	Type      int        `gorm:"default:0;not null;column:f_type;comment:'类型 1 股票 2外汇 3加密货币'" json:"type"`
	CreatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP;column:f_created_at;comment:'创建时间'" json:"created_at"`
	UpdatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP;column:f_updated_at;onUpdate:CURRENT_TIMESTAMP;comment:'修改时间'" json:"updated_at"`
}

func (AlphaInfo) TableName() string {
	return "t_alpha_info"
}

type AlphaInfoRepo struct {
	db  *gorm.DB
	ctx context.Context
}

func NewAlphaInfoRepo(ctx context.Context) *AlphaInfoRepo {
	return &AlphaInfoRepo{
		db:  public.DB.WithContext(ctx),
		ctx: ctx,
	}
}

func (a *AlphaInfoRepo) CreateAlphaInfo(alphaInfo *AlphaInfo) error {
	err := a.db.Create(alphaInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *AlphaInfoRepo) GetSymbolByType(t int) ([]*AlphaInfo, error) {
	var alphaInfos []*AlphaInfo
	err := a.db.Model(&AlphaInfo{}).Where("f_type = ?", t).Find(&alphaInfos).Error
	if err != nil {
		return nil, err
	}
	return alphaInfos, nil
}
