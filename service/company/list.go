package company

import (
	"financia/public/db/dao"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ListCompanyReq struct {
	Search   string   `form:"search"`
	Province []string `form:"province"`
	Page     int      `form:"page" binding:"required"`
	PageSize int      `form:"pageSize" binding:"required"`
}

type ListCompanyResp struct {
	List         []*ListCompanySimple `json:"list"`
	TotalPageNum int                  `json:"totalPageNum"`
	HasMore      bool                 `json:"hasMore"`
}

type ListCompanySimple struct {
	Id         int     `json:"id"`
	ComName    string  `json:"comName"`
	ComId      string  `json:"comId"`
	Chairman   string  `json:"chairman"`
	Manager    string  `json:"manager"`
	Secretary  string  `json:"secretary"`
	RegCapital float64 `json:"regCapital"`
	Province   string  `json:"province"`
	City       string  `json:"city"`
	Employees  int     `json:"employees"`
}

func ListCompany(c *gin.Context) {
	var req ListCompanyReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[ListCompany] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	zap.S().Debugf("[ListCompany] [req] = %#v", req)

	list, count, err := dao.GetCompanyList(c, req.Search, req.Province, req.Page, req.PageSize)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[ListCompany] [GetCompanyList] [err] = ", err.Error())
		return
	}

	respList := make([]*ListCompanySimple, 0, len(list))
	for _, v := range list {
		respList = append(respList, &ListCompanySimple{
			Id:         v.ID,
			ComName:    v.ComName,
			ComId:      v.ComID,
			Chairman:   v.Chairman,
			Manager:    v.Manager,
			Secretary:  v.Secretary,
			RegCapital: v.RegCapital,
			Province:   v.Province,
			City:       v.City,
			Employees:  v.Employees,
		})
	}

	util.SuccessResp(c, &ListCompanyResp{
		List:         respList,
		HasMore:      count > int64(req.Page*(req.PageSize-1)+len(list)),
		TotalPageNum: int(count/int64(req.PageSize) + 1),
	})
}
