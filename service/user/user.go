package user

import (
	"errors"
	"financia/public"
	"financia/public/db/dao"
	"financia/server"
	"financia/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// Code 获取验证码
func Code(c *gin.Context) {
	var req GetCodeReq
	if err := c.ShouldBindQuery(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[Code] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	code, err := dao.GetEmailCode(c, req.Email)
	if err == nil || code != "" {
		util.FailRespWithCode(c, util.CodeLimitError)
		zap.S().Error("[Code] [GetEmailCode] [err] = 验证码发送过于频繁")
		return
	}

	code = public.GenerateVerificationCode(6)

	go server.SendEmail(req.Email, public.EmailTitle, code)

	err = dao.SetEmailCode(c, req.Email, code)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[Code] [SetEmailCode] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, nil)
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

func Info(c *gin.Context) {
	userId := util.GetUid(c)
	user, err := dao.GetUser(c, userId)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[Info] [GetUserInfo] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, &UserInfoResp{
		Email:    user.Email,
		UserName: user.Username,
	})
}
