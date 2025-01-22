package tushare

type DailyReq struct {
	TsCode     string `json:"ts_code,omitempty"`
	TradeDate  string `json:"trade_date,omitempty"`
	StartDate  string `json:"start_date,omitempty"`
	Exchange   string `json:"exchange,omitempty"`
	EndDate    string `json:"end_date,omitempty"`
	Prd        string `json:"prd,omitempty"`
	ReportType int    `json:"report_type,omitempty"` // 1
	Q          string `json:"q,omitempty"`
	StartM     string `json:"start_m,omitempty"`
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

type StockIncomeResp struct {
	AnnDate      string  `json:"annDate"`      // 公告日期
	BasicEps     float64 `json:"basicEps"`     // 基本每股收益
	TotalRevenue float64 `json:"totalRevenue"` // 营业总收入
	TotalCogs    float64 `json:"totalCogs"`    // 营业总成本
	OperExp      float64 `json:"operExp"`      // 营业支出
	TotalProfit  float64 `json:"totalProfit"`  // 利润总额
	IncomeTax    float64 `json:"incomeTax"`    // 所得税费用
	NIncome      float64 `json:"nIncome"`      // 净利润
	TComprIncome float64 `json:"tComprIncome"` // 综合收益总额
}

type StockForecastResp struct {
	AnnDate       string  `json:"annDate"`       // 公告日期
	Type          string  `json:"type"`          // 预告类型
	PChangeMin    float64 `json:"pChangeMin"`    // 预告净利润变动幅度下限
	PChangeMax    float64 `json:"pChangeMax"`    // 预告净利润变动幅度上限
	NetProfitMin  float64 `json:"netProfitMin"`  // 预告净利润下限
	NetProfitMax  float64 `json:"netProfitMax"`  // 预告净利润上限
	LastParentNet float64 `json:"lastParentNet"` // 上年同期归属母公司净利润
	ChangeReason  string  `json:"changeReason"`  // 预告净利润变动原因
}

type StockTop10Resp struct {
	AnnDate        string  `json:"annDate"`        // 公告日期
	HolderName     string  `json:"holderName"`     // 股东名称
	HoldAmount     float64 `json:"holdAmount"`     // 持股数量
	HoldRatio      float64 `json:"holdRatio"`      // 持股比例
	HoldFloatRatio float64 `json:"holdFloatRatio"` // 流通股比例
	HoldChange     float64 `json:"holdChange"`     // 变动
	HolderType     string  `json:"holderType"`     // 股东类型
}

type StockHsgtTop10Resp struct {
	Name   string  `json:"name"`   // 股票名称
	Close  float64 `json:"close"`  // 收盘价
	Change float64 `json:"change"` // 涨跌幅
	Rank   int     `json:"rank"`   // 排名
	Amount float64 `json:"amount"` // 持股数量
}

type EconomicsShiborResp struct {
	Date   string  `json:"date"`
	On     float64 `json:"on"`
	OneW   float64 `json:"oneW"`
	TwoW   float64 `json:"twoW"`
	OneM   float64 `json:"oneM"`
	ThreeM float64 `json:"threeM"`
	SixM   float64 `json:"sixM"`
	NineM  float64 `json:"nineM"`
	OneY   float64 `json:"oneY"`
}

type EconomicsCnGDPResp struct {
	GDP    float64 `json:"gdp"`
	GDPYoy float64 `json:"gdpYoy"`
	PI     float64 `json:"pi"`
	PIYoy  float64 `json:"piYoy"`
	SI     float64 `json:"si"`
	SIYoy  float64 `json:"siYoy"`
	TI     float64 `json:"ti"`
	TIYoy  float64 `json:"tiYoy"`
}

type EconomicsCnCPIResp struct {
	Month    string  `json:"month"`
	NtYoy    float64 `json:"ntYoy"`
	NtMom    float64 `json:"ntMom"`
	NtAccu   float64 `json:"ntAccu"`
	TownYoy  float64 `json:"townYoy"`
	TownMom  float64 `json:"townMom"`
	TownAccu float64 `json:"townAccu"`
	CntYoy   float64 `json:"cntYoy"`
	CntMom   float64 `json:"cntMom"`
	CntAccu  float64 `json:"cntAccu"`
}
