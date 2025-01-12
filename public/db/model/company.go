package model

type CompanyInfo struct {
	ID            int     `gorm:"column:f_id;primaryKey;autoIncrement;comment:'主键'"`
	TsCode        string  `gorm:"column:f_ts_code;type:varchar(20);not null;comment:'股票代码'"`
	ComName       string  `gorm:"column:f_com_name;type:varchar(200);default:'';comment:'公司全称'"`
	ComID         string  `gorm:"column:f_com_id;type:varchar(50);default:'';comment:'统一社会信用代码'"`
	Chairman      string  `gorm:"column:f_chairman;type:varchar(100);default:'';comment:'法人代表'"`
	Manager       string  `gorm:"column:f_manager;type:varchar(100);default:'';comment:'总经理'"`
	Secretary     string  `gorm:"column:f_secretary;type:varchar(100);default:'';comment:'董秘'"`
	RegCapital    float64 `gorm:"column:f_reg_capital;type:float;default:0;not null;comment:'注册资本'"`
	Province      string  `gorm:"column:f_province;type:varchar(50);default:'';comment:'所在省份'"`
	City          string  `gorm:"column:f_city;type:varchar(50);default:'';not null;comment:'所在城市'"`
	Introduction  string  `gorm:"column:f_introduction;type:text;comment:'公司介绍'"`
	BusinessScope string  `gorm:"column:f_business_scope;type:text;comment:'经营范围'"`
	Employees     int     `gorm:"column:f_employees;type:int;default:0;comment:'员工人数'"`
	MainBusiness  string  `gorm:"column:f_main_business;type:text;comment:'主要业务及产品'"`
}

// TableName specifies the table name for the TCompanyInfo model.
func (CompanyInfo) TableName() string {
	return "t_company_info"
}
