package stock

import (
	"financia/public/db/dao"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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

func ListStock(c *gin.Context) {
	var req ListStockReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[ListStock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	zap.S().Debugf("[ListStock] [req] = %#v", req)

	list, count, err := dao.GetStockList(c, req.Search, req.IsHs, req.Exchange, req.Market, req.Page, req.PageSize)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[ListStock] [GetStockList] [err] = ", err.Error())
		return
	}

	respList := make([]*ListStockSimple, 0, len(list))
	for _, v := range list {
		respList = append(respList, &ListStockSimple{
			Id:         v.Id,
			Name:       v.Name,
			Area:       v.Area,
			Industry:   v.Industry,
			Market:     v.Market,
			ActName:    v.ActName,
			ActEntType: v.ActEntType,
			FullName:   v.FullName,
			EnName:     v.EnName,
			CnSpell:    v.CnSpell,
			Exchange:   v.Exchange,
			CurrType:   v.CurrType,
			ListStatus: v.ListStatus,
			IsHs:       v.IsHs,
		})
	}

	util.SuccessResp(c, &ListStockResp{
		List:         respList,
		HasMore:      count > int64(req.Page*(req.PageSize-1)+len(list)),
		TotalPageNum: int(count/int64(req.PageSize) + 1),
	})
}
