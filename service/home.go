package service

import (
	"financia/models"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetInfoReq struct {
	Type int `form:"type" binding:"required"`
}

type GetInfoRespSimple struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type GetInfoResp struct {
	List []*GetInfoRespSimple `json:"list"`
}

func GetInfo(c *gin.Context) {
	var query GetInfoReq
	err := c.ShouldBindQuery(&query)
	if err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[GetInfo] 参数错误")
		return
	}

	alphaInfoRepo := models.NewAlphaInfoRepo(c)
	infoList, err := alphaInfoRepo.GetSymbolByType(query.Type)
	if err != nil {
		util.FailResp(c, err.Error())
		return
	}

	resp := &GetInfoResp{
		List: make([]*GetInfoRespSimple, 0, len(infoList)),
	}
	for _, info := range infoList {
		resp.List = append(resp.List, &GetInfoRespSimple{
			Name:   info.Name,
			Symbol: info.Symbol,
		})
	}

	util.SuccessResp(c, resp)
}
