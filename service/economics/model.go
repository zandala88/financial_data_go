package economics

import "financia/server/tushare"

type ShiborEconomicsResp struct {
	List []*tushare.EconomicsShiborResp `json:"list"`
}

type CnGdpEconomicsReq struct {
	Year    string `form:"year" binding:"required"`
	Quarter string `form:"quarter" binding:"required"`
}

type CnGdpEconomicsResp struct {
	List []*tushare.EconomicsCnGDPResp `json:"list"`
}
