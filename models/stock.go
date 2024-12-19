package models

type Stock struct {
	Id      int64   `gorm:"primaryKey;autoIncrement;column:f_id;comment:'主键 ID，自增'" json:"id"`
	Company string  `gorm:"not null;column:f_company;comment:'公司名称'" json:"company"`
	Date    string  `gorm:"type:date;not null;column:f_date;comment:'日期'" json:"date"`
	Open    float64 `gorm:"type:decimal(10,3);default:0.000;not null;column:f_open;comment:'开盘价'" json:"open"`
	High    float64 `gorm:"type:decimal(10,3);default:0.000;not null;column:f_high;comment:'最高价'" json:"high"`
	Low     float64 `gorm:"type:decimal(10,3);default:0.000;not null;column:f_low;comment:'最低价'" json:"low"`
	Close   float64 `gorm:"type:decimal(10,3);default:0.000;not null;column:f_close;comment:'收盘价'" json:"close"`
	Volume  int64   `gorm:"default:0;not null;column:f_volume;comment:'成交量'" json:"volume"`
}

func (*Stock) TableName() string {
	return "t_stock_data"
}
