package tushare

type DailyReq struct {
	TsCode    string `json:"ts_code,omitempty"`
	TradeDate string `json:"trade_date,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
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
