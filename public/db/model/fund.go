package model

import "time"

type FundInfo struct {
	Id           int64   `gorm:"column:id;primaryKey;autoIncrement;comment:'主键'"`                           // 主键
	TsCode       string  `gorm:"column:f_ts_code;type:varchar(255);default:'';not null;comment:'TS代码'"`     // TS代码
	Name         string  `gorm:"column:f_name;type:varchar(255);default:'';not null;comment:'简称'"`          // 简称
	Management   string  `gorm:"column:f_management;type:varchar(255);default:'';not null;comment:'管理人'"`   // 管理人
	Custodian    string  `gorm:"column:f_custodian;type:varchar(255);default:'';not null;comment:'托管人'"`    // 托管人
	FundType     string  `gorm:"column:f_fund_type;type:varchar(255);default:'';not null;comment:'投资类型'"`   // 投资类型
	IssueAmount  float64 `gorm:"column:f_issue_amount;default:0;comment:'发行份额(亿份)'"`                        // 发行份额(亿份)
	MFee         float64 `gorm:"column:f_m_fee;default:0;comment:'管理费'"`                                    // 管理费
	CCFee        float64 `gorm:"column:f_c_fee;default:0;comment:'托管费'"`                                    // 托管费
	DurationYear float64 `gorm:"column:f_duration_year;default:0;comment:'存续期'"`                            // 存续期
	PValue       float64 `gorm:"column:f_p_value;default:0;comment:'面值'"`                                   // 面值
	MinAmount    float64 `gorm:"column:f_min_amount;default:0;comment:'起点金额(万元)'"`                          // 起点金额(万元)
	ExpReturn    float64 `gorm:"column:f_exp_return;default:0;comment:'预期收益率'"`                             // 预期收益率
	Benchmark    string  `gorm:"column:f_benchmark;type:varchar(255);default:'';not null;comment:'业绩比较基准'"` // 业绩比较基准
	InvestType   string  `gorm:"column:f_invest_type;type:varchar(255);default:'';not null;comment:'投资风格'"` // 投资风格
	Type         string  `gorm:"column:f_type;type:varchar(255);default:'';not null;comment:'基金类型'"`        // 基金类型
	Trustee      string  `gorm:"column:f_trustee;type:varchar(255);default:'';not null;comment:'受托人'"`      // 受托人
	Flag         int     `gorm:"column:f_flag;default:0;comment:'标记'"`                                      // 标记
}

func (FundInfo) TableName() string {
	return "t_fund_info"
}

type FundData struct {
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
	Vol       float64   `gorm:"type:bigint;column:f_vol" json:"vol"`
	Amount    float64   `gorm:"type:decimal(20,2);column:f_amount" json:"amount"`
}

func (FundData) TableName() string {
	return "t_fund_data"
}
