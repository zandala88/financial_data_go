package fund

import "financia/server/tushare"

type DataFundReq struct {
	Id        int    `form:"id" binding:"required"`
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
}

type DataFundResp struct {
	Have bool              `json:"have"`
	List []*DataFundSimple `json:"list"`
}

type DataFundSimple struct {
	TradeDate string  `json:"tradeDate"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	PreClose  float64 `json:"preClose"`
	Change    float64 `json:"change"`
	PctChg    float64 `json:"pctChg"`
	Vol       int64   `json:"vol"`
	Amount    float64 `json:"amount"`
}

type GraphFundResp struct {
	Radio []*tushare.FundSalesRatioResp `json:"radio"`
	Inst  []*tushare.FundSalesVolResp   `json:"inst"`
}

type HaveFundReq struct {
	Id int `form:"id" binding:"required"`
}

type HaveFundResp struct {
	Have bool `json:"have"`
}

type ListFundReq struct {
	Search     string   `form:"search"`
	FundType   []string `form:"fundType"`
	InvestType []string `form:"investType"`
	Page       int      `form:"page" binding:"required"`
	PageSize   int      `form:"pageSize" binding:"required"`
}

type ListFundResp struct {
	List         []*ListFundSimple `json:"list"`
	TotalPageNum int               `json:"totalPageNum"`
	HasMore      bool              `json:"hasMore"`
}

type ListFundSimple struct {
	Id           int64   `json:"id"`
	Name         string  `json:"name"`
	Management   string  `json:"management"`
	Custodian    string  `json:"custodian"`
	FundType     string  `json:"fundType"`
	IssueAmount  float64 `json:"issueAmount"`
	MFree        float64 `json:"mFree"`
	CFree        float64 `json:"cFree"`
	DurationYear float64 `json:"durationYear"`
	PValue       float64 `json:"pValue"`
	MinAmount    float64 `json:"minAmount"`
	ExpReturn    float64 `json:"expReturn"`
	Benchmark    string  `json:"benchmark"`
	InvestType   string  `json:"investType"`
	Type         string  `json:"type"`
	Trustee      string  `json:"trustee"`
}

type QueryFundResp struct {
	FundTypeList   []string `json:"fundTypeList"`
	InvestTypeList []string `json:"investTypeList"`
}
