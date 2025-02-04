package stock

import "financia/server/tushare"

type DataStockReq struct {
	Id        int    `form:"id" binding:"required"`
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
}

type DataStockResp struct {
	Have bool               `json:"have"`
	List []*DataStockSimple `json:"list"`
}

type DataStockSimple struct {
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

type GraphStockResp struct {
	IsHs     map[string]int `json:"isHs"`
	Exchange map[string]int `json:"exchange"`
	Market   map[string]int `json:"market"`
}

type HaveStockReq struct {
	Id int `form:"id" binding:"required"`
}

type HaveStockResp struct {
	Have bool `json:"have"`
}

type IncomeStockReq struct {
	Id int `form:"id" binding:"required"`
}

type IncomeStockResp struct {
	List []*tushare.StockIncomeResp `json:"list"`
}

type InfoStockReq struct {
	Id int `form:"id" binding:"required"`
}

type InfoStockResp struct {
	FullName string `json:"name"`
	Industry string `json:"industry"`
	Market   string `json:"market"`
}

type ListStockReq struct {
	Search   string   `form:"search"`
	IsHs     []string `form:"isHs"`
	Exchange []string `form:"exchange"`
	Market   []string `form:"market"`
	Page     int      `form:"page" binding:"required"`
	PageSize int      `form:"pageSize" binding:"required"`
}

type ListStockResp struct {
	List         []*ListStockSimple `json:"list"`
	TotalPageNum int                `json:"totalPageNum"`
	HasMore      bool               `json:"hasMore"`
}

type ListStockSimple struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Area       string `json:"area"`
	Industry   string `json:"industry"`
	Market     string `json:"market"`
	ActName    string `json:"actName"`
	ActEntType string `json:"actEntType"`
	FullName   string `json:"fullName"`
	EnName     string `json:"enName"`
	CnSpell    string `json:"cnSpell"`
	Exchange   string `json:"exchange"`
	CurrType   string `json:"currType"`
	ListStatus string `json:"listStatus"`
	IsHs       string `json:"isHs"`
}

type QueryStockResp struct {
	IsHsList     []string `json:"isHsList"`
	ExchangeList []string `json:"exchangeList"`
	MarketList   []string `json:"marketList"`
}

type ForecastStockReq struct {
	Id int `form:"id" binding:"required"`
}

type ForecastStockResp struct {
	List []*tushare.StockForecastResp `json:"list"`
}

type Top10StockReq struct {
	Id int `form:"id" binding:"required"`
}

type Top10StockResp struct {
	List []*tushare.StockTop10Resp `json:"list"`
	Rank []*Top10StockRank         `json:"rank"`
}

type Top10StockRank struct {
	HolderName string  `json:"holderName"`
	HoldRatio  float64 `json:"holdRatio"`
}

type Top10HsgtStockResp struct {
	ShList []*tushare.StockHsgtTop10Resp `json:"shList"`
	SzList []*tushare.StockHsgtTop10Resp `json:"szList"`
}

type PredictStockReq struct {
	Id int `form:"id" binding:"required"`
}

type PredictStockResp struct {
	List []float64 `json:"list"`
	Val  float64   `json:"val"`
}

type PythonPredictReq struct {
	Data []*PythonPredictReqSimple `json:"data"`
}

type PythonPredictReqSimple struct {
	Date   string  `json:"date"`
	CoIMF1 float64 `json:"Co-IMF1"`
	CoIMF2 float64 `json:"Co-IMF2"`
	CoIMF3 float64 `json:"Co-IMF3"`
	CoIMF4 int64   `json:"Co-IMF4"`
	Target float64 `json:"Target"`
}

type PythonPredictResp struct {
	Code int                   `json:"code"`
	Data PythonPredictRespData `json:"data"`
}

type PythonPredictRespData struct {
	Val float64 `json:"val"`
}
