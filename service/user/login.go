package user

import (
	"financia/public/db/dao"
	"financia/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoginReq struct {
	Email    string `form:"email"  binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type LoginResp struct {
	Token string `json:"token"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[Login] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	userId := dao.GetUserId(c, req.Email)
	if userId <= 0 {
		util.FailRespWithCode(c, util.ReqDataError)
		zap.S().Error("[Login] [GetUserId] [err] = 用户不存在")
		return
	}

	token, err := util.GenerateJWT(userId)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[Login] [GenerateJWT] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, LoginResp{Token: token})
}
