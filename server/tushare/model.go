package tushare

type DailyReq struct {
	TsCode    string `json:"ts_code,omitempty"`
	TradeDate string `json:"trade_date,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	Exchange  string `json:"exchange,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Prd       string `json:"prd,omitempty"`
}

type DailyResp struct {
	Fields []string        `json:"fields"`
	Items  [][]interface{} `json:"items"`
}

type FundSalesRatioResp struct {
	Year      string  `json:"year"`
	Bank      float64 `json:"bank"`
	SecComp   float64 `json:"secComp"`
	FundComp  float64 `json:"fundComp"`
	IndepComp float64 `json:"indepComp"`
	Rests     float64 `json:"rests"`
}

type FundSalesVolResp struct {
	Year      string  `json:"year"`
	Quarter   string  `json:"quarter"`
	InstName  string  `json:"instName"`
	FundScale float64 `json:"fundScale"`
	Scale     float64 `json:"scale"`
	Rank      int     `json:"rank"`
}

type FutTradeCalResp struct {
	CalDate string `json:"calDate"`
	IsOpen  int    `json:"isOpen"` // 0: 休市 1: 开市
}

type FutWeeklyDetailResp struct {
	Vol          int     `json:"vol"`
	VolYoy       float64 `json:"volYoy"`
	Amount       float64 `json:"amount"`
	AmountYoy    float64 `json:"amountYoy"`
	CumVol       int     `json:"cumVol"`
	CumVolYoy    float64 `json:"cumVolYoy"`
	Cumamt       float64 `json:"cumamt"`
	CumamtYoy    float64 `json:"cumamtYoy"`
	OpenInterest int     `json:"openInterest"`
	InterestWow  float64 `json:"interestWow"`
	McClose      float64 `json:"mcClose"`
	CloseWow     float64 `json:"closeWow"`
	WeekDate     string  `json:"weekDate"`
}
