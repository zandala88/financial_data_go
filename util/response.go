package util

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SuccessResp(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

func FailRespWithCode(c *gin.Context, code int) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  errMsg[code],
	})
}

func FailRespWithCodeAndZap(c *gin.Context, code int, format, logStr string) {
	FailRespWithCode(c, code)
	zap.S().Errorf(format, logStr)
}

var errMsg = map[int]string{
	Ok: "Success",

	InternalServerError: "服务器内部错误",
	InvalidToken:        "无效的Token",
	ShouldBindJSONError: "参数错误",
	ReqDataError:        "参数内容错误",
	CodeLimitError:      "验证码发送过于频繁",
}

const (
	Ok           = 200
	InvalidToken = 403

	InternalServerError = iota + 2001

	ShouldBindJSONError
	ReqDataError
	CodeLimitError
)
