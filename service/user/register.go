package user

import (
	"errors"
	"financia/public/db/dao"
	"financia/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RegisterReq struct {
	Email    string `form:"email"  binding:"required,email"`
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Code     string `form:"code" binding:"required"`
}

type RegisterResp struct {
	Token string `json:"token"`
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[Register] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	code, err := dao.GetEmailCode(c, req.Email)
	if err != nil && !errors.Is(err, redis.Nil) {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[Register] [GetEmailCode] [err] = ", err)
		return
	}
	if code != req.Code {
		util.FailRespWithCode(c, util.ReqDataError)
		zap.S().Errorf("[Register] 验证码错误 req.Code = %s code = %s", req.Code, code)
	}

	userId := dao.GetUserId(c, req.Email)
	if userId > 0 {
		util.FailRespWithCode(c, util.ReqDataError)
		zap.S().Error("[Register] [GetUserId] [err] = 用户已存在")
		return
	}

	if err := dao.CreateUser(c, req.Email, req.Username, req.Password); err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[Register] [CreateUser] [err] = ", err.Error())
		return
	}

	userId = dao.GetUserId(c, req.Username)
	token, err := util.GenerateJWT(userId)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[Register] [GenerateJWT] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, RegisterResp{Token: token})
}
