package fut

import "financia/server/tushare"

// CU 铜
// SR 白糖
// CF 棉花
// AL 铝
// ZN 锌
// JD 鸡蛋
// FG 玻璃
// AP 苹果
// PP 聚丙烯
// RB 螺纹钢
type QueryFutResp struct {
	List []*QueryFutSimple `json:"list"`
}

type QueryFutSimple struct {
	Prd  string `json:"prd"`
	Name string `json:"name"`
}

type DetailFutReq struct {
	Prd string `form:"prd" binding:"required"`
}

type DetailFutResp struct {
	List []*tushare.FutWeeklyDetailResp `json:"list"`
}

type CalFutResp struct {
	Sse  []*tushare.FutTradeCalResp `json:"sse"`
	Szse []*tushare.FutTradeCalResp `json:"szse"`
}
