package company

import (
	"financia/public/db/dao"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DetailCompanyReq struct {
	Id int `form:"id" binding:"required"`
}

type DetailCompanyResp struct {
	ComName       string  `json:"comName"`
	ComId         string  `json:"comId"`
	Chairman      string  `json:"chairman"`
	Manager       string  `json:"manager"`
	Secretary     string  `json:"secretary"`
	RegCapital    float64 `json:"regCapital"`
	Province      string  `json:"province"`
	City          string  `json:"city"`
	Employees     int     `json:"employees"`
	Introduction  string  `json:"introduction"`
	BusinessScope string  `json:"businessScope"`
	MainBusiness  string  `json:"mainBusiness"`
}

func DetailCompany(c *gin.Context) {
	var req DetailCompanyReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[DetailCompany] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	company, err := dao.GetCompany(c, req.Id)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[DetailCompany] [GetCompanyDetail] [err] = ", err.Error())
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
