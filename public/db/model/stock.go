package model

import "time"

type StockInfo struct {
	Id         int    `gorm:"column:f_id;primaryKey;autoIncrement" json:"id"`
	TsCode     string `gorm:"column:f_ts_code;type:varchar(20);not null" json:"tsCode"`
	Symbol     string `gorm:"column:f_symbol;type:varchar(10);not null" json:"symbol"`
	Name       string `gorm:"column:f_name;type:varchar(100);not null" json:"name"`
	Area       string `gorm:"column:f_area;type:varchar(50)" json:"area"`
	Industry   string `gorm:"column:f_industry;type:varchar(50)" json:"industry"`
	Market     string `gorm:"column:f_market;type:varchar(50)" json:"market"`
	ActName    string `gorm:"column:f_act_name;type:varchar(100)" json:"actName"`
	ActEntType string `gorm:"column:f_act_ent_type;type:varchar(50)" json:"actEntType"`
	FullName   string `gorm:"column:f_fullname;type:varchar(200)" json:"fullName"`
	EnName     string `gorm:"column:f_enname;type:varchar(200)" json:"enName"`
	CnSpell    string `gorm:"column:f_cnspell;type:varchar(50)" json:"cnSpell"`
	Exchange   string `gorm:"column:f_exchange;type:varchar(20)" json:"exchange"`
	CurrType   string `gorm:"column:f_curr_type;type:varchar(10)" json:"currType"`
	ListStatus string `gorm:"column:f_list_status;type:enum('L','D','P')" json:"listStatus"`
	IsHs       string `gorm:"column:f_is_hs;type:enum('N','H','S')" json:"isHs"`
}

func (StockInfo) TableName() string {
	return "t_stock_info"
}

type StockData struct {
	Id        int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TsCode    string    `gorm:"type:varchar(20);column:f_ts_code" json:"tsCode"`
	TradeDate time.Time `gorm:"type:date;column:f_trade_date" json:"tradeDate"`
	Open      float64   `gorm:"type:decimal(10,2);column:f_open" json:"open"`
	High      float64   `gorm:"type:decimal(10,2);column:f_high" json:"high"`
	Low       float64   `gorm:"type:decimal(10,2);column:f_low" json:"low"`
	Close     float64   `gorm:"type:decimal(10,2);column:f_close" json:"close"`
	PreClose  float64   `gorm:"type:decimal(10,2);column:f_pre_close" json:"preClose"`
	Change    float64   `gorm:"type:decimal(10,2);column:f_change" json:"change"`
	PctChg    float64   `gorm:"type:decimal(5,2);column:f_pct_chg" json:"pctChg"`
	Vol       int64     `gorm:"type:bigint;column:f_vol" json:"vol"`
	Amount    float64   `gorm:"type:decimal(20,2);column:f_amount" json:"amount"`
}

func (StockData) TableName() string {
	return "t_stock_data"
}

type StockPredict struct {
	Id        int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TsCode    string    `gorm:"type:varchar(20);column:f_ts_code" json:"tsCode"`
	TradeDate time.Time `gorm:"type:date;column:f_trade_date" json:"tradeDate"`
	Predict   float64   `gorm:"type:decimal(10,2);column:f_predict" json:"predict"`
}

func (StockPredict) TableName() string {
	return "t_stock_predict"
}
