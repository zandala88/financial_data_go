package user

import (
	"financia/public"
	"financia/public/db/dao"
	"financia/server"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetCodeReq struct {
	Email string `form:"email" binding:"required,email"`
}

type GetCodeResp struct {
	Code string `json:"code"`
}

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
