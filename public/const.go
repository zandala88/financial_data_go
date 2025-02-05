package public

const (
	RedisKeyFundSalesRatio  = "fund_sales_ratio"
	RedisKeyShiborEconomics = "shibor_economics"
	RedisKeyCnGdpEconomics  = "cn_gdp_economics"
	RedisKeyCnCpiEconomics  = "cn_cpi_economics"
	RedisKeyFundRadio       = "fund_radio"
	RedisKeyFundVol         = "fund_vol"
	RedisKeyGraphStock      = "graph_stock"
	RedisKeyStockPredict    = "stock_predict:%d"
	RedisKeyStockFollow     = "stock_follow:%d"
	RedisKeyFundFollow      = "fund_follow:%d"
)

const (
	TuShareDaily            = "daily"
	TuShareFundDaily        = "fund_daily"
	TuShareFundSalesRatio   = "fund_sales_ratio"
	TuShareFundSalesVol     = "fund_sales_vol"
	TuShareFutTradeCal      = "trade_cal"
	TuShareFutWeeklyDetail  = "fut_weekly_detail"
	TuShareStockIncome      = "income"
	TuShareStockForecast    = "forecast"
	TuShareStockHolderTop10 = "top10_holders"
	TuShareStockHsgtTop10   = "hsgt_top10"
	TuShareEconomicsShibor  = "shibor"
	TuShareEconomicsCnGDP   = "cn_gdp"
	TuShareEconomicsCnCPI   = "cn_cpi"
)

const (
	EmailTitle = "zandala-financial 验证码"
)

const (
	FundInfoFlagExist = 1
)

const (
	CnGdpEconomicsStartYear  = "1992"
	CnGdpEconomicsEndYear    = "2024"
	CnGdpEconomicsEndQuarter = "4"
)

// 0: 休市 1: 开市
const (
	MarketStatusClose = iota
	MarketStatusOpen
)
