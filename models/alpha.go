package models

import "time"

type AlphaInfo struct {
	Id        int64      `gorm:"primaryKey;autoIncrement;column:f_id;comment:'主键 ID'" json:"id"`
	Name      string     `gorm:"default:'';not null;column:f_name;comment:'名称'" json:"name"`
	Symbol    string     `gorm:"default:'';not null;column:f_symbol;comment:'Alpha Vantage API 参数'" json:"symbol"`
	SymbolTo  string     `gorm:"default:'';not null;column:f_symbol_to;comment:'外汇参数'" json:"symbol_to"`
	Type      int        `gorm:"default:0;not null;column:f_type;comment:'类型，区分股票或外汇（使用整数表示，默认0）'" json:"type"`
	CreatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP;column:f_created_at;comment:'创建时间'" json:"created_at"`
	UpdatedAt *time.Time `gorm:"default:CURRENT_TIMESTAMP;column:f_updated_at;onUpdate:CURRENT_TIMESTAMP;comment:'修改时间'" json:"updated_at"`
}

func (AlphaInfo) TableName() string {
	return "t_alpha_info"
}