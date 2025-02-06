package company

import (
	"financia/public/db/dao"
	"financia/util"
	"github.com/gin-gonic/gin"
)

func DetailCompany(c *gin.Context) {
	var req DetailCompanyReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[DetailCompany] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	company, err := dao.GetCompany(c, req.Id)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[DetailCompany] [GetCompany] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, &DetailCompanyResp{
		ComName:       company.ComName,
		ComId:         company.ComID,
		Chairman:      company.Chairman,
		Manager:       company.Manager,
		Secretary:     company.Secretary,
		RegCapital:    company.RegCapital,
		Province:      company.Province,
		City:          company.City,
		Employees:     company.Employees,
		Introduction:  company.Introduction,
		BusinessScope: company.BusinessScope,
		MainBusiness:  company.MainBusiness,
	})
}

func ListCompany(c *gin.Context) {
	var req ListCompanyReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[ListCompany] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	list, count, err := dao.GetCompanyList(c, req.Search, req.Province, req.Page, req.PageSize)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[ListCompany] [GetCompanyList] [err] = ", err.Error())
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

func QueryCompany(c *gin.Context) {
	dis, err := dao.ProvinceDis(c)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[QueryCompany] [ProvinceDis] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, &QueryCompanyResp{
		List: dis,
	})
}
