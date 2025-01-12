package company

import (
	"financia/public/db/dao"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type QueryCompanyResp struct {
	List []string `json:"list"`
}

func QueryCompany(c *gin.Context) {
	dis, err := dao.ProvinceDis(c)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[QueryCompany] [ProvinceDis] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, &QueryCompanyResp{
		List: dis,
	})
}
