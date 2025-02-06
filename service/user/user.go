package user

import (
	"errors"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/dao"
	"financia/server"
	"financia/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math"
	"time"
)

// Code 获取验证码
func Code(c *gin.Context) {
	var req GetCodeReq
	if err := c.ShouldBindQuery(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[Code] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	code, err := dao.GetEmailCode(c, req.Email)
	if err == nil || code != "" {
		util.FailRespWithCodeAndZap(c, util.CodeLimitError, "[Code] [GetEmailCode] [err] = ", "验证码发送过于频繁")
		return
	}

	code = public.GenerateVerificationCode(6)

	go server.SendEmail(req.Email, public.EmailTitle, code)

	err = dao.SetEmailCode(c, req.Email, code)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Code] [SetEmailCode] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, nil)
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[Login] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	userId := dao.GetUserId(c, req.Email)
	if userId <= 0 {
		util.FailRespWithCodeAndZap(c, util.ReqDataError, "[Login] [GetUserId] [err] = ", "用户不存在")
		return
	}

	token, err := util.GenerateJWT(userId)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Login] [GenerateJWT] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, LoginResp{Token: token})
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[Register] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	code, err := dao.GetEmailCode(c, req.Email)
	if err != nil && !errors.Is(err, redis.Nil) {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Register] [GetEmailCode] [err] = ", err.Error())
		return
	}
	if code != req.Code {
		util.FailRespWithCodeAndZap(c, util.ReqDataError, "[Register] [GetEmailCode] [err] = ", "验证码错误")
		return
	}

	userId := dao.GetUserId(c, req.Email)
	if userId > 0 {
		util.FailRespWithCodeAndZap(c, util.ReqDataError, "[Register] [GetUserId] [err] = ", "用户已存在")
		return
	}

	if err := dao.CreateUser(c, req.Email, req.Username, req.Password); err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Register] [CreateUser] [err] = ", err.Error())
		return
	}

	userId = dao.GetUserId(c, req.Username)
	token, err := util.GenerateJWT(userId)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Register] [GenerateJWT] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, RegisterResp{Token: token})
}

func Info(c *gin.Context) {
	userId := util.GetUid(c)
	user, err := dao.GetUser(c, userId)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Info] [GetUser] [err] = ", err.Error())
		return
	}

	resp := &UserInfoResp{
		Email:     user.Email,
		UserName:  user.Username,
		StockList: make([]*UserInfoData, 0),
		FundList:  make([]*UserInfoData, 0),
	}

	if err := predict(c, userId, resp); err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Info] [predict] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, resp)
}

func Tip(c *gin.Context) {
	// 是否登陆了
	userId := util.GetUid(c)
	if userId == public.EmptyUserId {
		util.SuccessResp(c, TipResp{Exists: false})
		return
	}

	// 是否已经不再提醒了
	key := fmt.Sprintf(public.RedisKeyTip, userId)
	rdb := connector.GetRedis().WithContext(c)
	exists, err := rdb.Exists(c, key).Result()
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Tip] [Exists] [err] = ", err.Error())
		return
	}
	if exists == public.RedisExists {
		util.SuccessResp(c, TipResp{Exists: false})
		return
	}

	// 获取提示信息
	pResp := &UserInfoResp{
		StockList: make([]*UserInfoData, 0),
		FundList:  make([]*UserInfoData, 0),
	}

	if err := predict(c, userId, pResp); err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Info] [predict] [err] = ", err.Error())
		return
	}

	var (
		stockRise TipSimple
		stockFall TipSimple
		fundRise  TipSimple
		fundFall  TipSimple
	)

	for _, v := range pResp.StockList {
		if v.NextVal == 0 {
			continue
		}
		if v.NextVal > v.Val {
			if stockRise.Val == 0 || stockRise.Scope < 100*(v.NextVal-v.Val)/v.Val {
				stockRise = TipSimple{
					Name:  v.Name,
					Val:   v.NextVal,
					Scope: math.Floor((v.NextVal-v.Val)/v.Val*10000) / 100,
				}
			}
		} else {
			if stockFall.Val == 0 || stockFall.Scope > 100*(v.NextVal-v.Val)/v.Val {
				stockFall = TipSimple{
					Name:  v.Name,
					Val:   v.NextVal,
					Scope: math.Floor((v.NextVal-v.Val)/v.Val*10000) / 100,
				}
			}
		}
	}

	for _, v := range pResp.FundList {
		if v.NextVal == 0 {
			continue
		}
		if v.NextVal > v.Val {
			if fundRise.Val == 0 || fundRise.Scope < 100*(v.NextVal-v.Val)/v.Val {
				fundRise = TipSimple{
					Name:  v.Name,
					Val:   v.NextVal,
					Scope: math.Floor((v.NextVal-v.Val)/v.Val*10000) / 100,
				}
			}
		} else {
			if fundFall.Val == 0 || fundFall.Scope > 100*(v.NextVal-v.Val)/v.Val {
				fundFall = TipSimple{
					Name:  v.Name,
					Val:   v.NextVal,
					Scope: math.Floor((v.NextVal-v.Val)/v.Val*10000) / 100,
				}
			}
		}
	}

	util.SuccessResp(c, TipResp{
		Exists:    true,
		StockRise: stockRise,
		StockFall: stockFall,
		FundRise:  fundRise,
		FundFall:  fundFall,
	})
}

func TipConfirm(c *gin.Context) {
	userId := util.GetUid(c)
	// 设置不再提醒
	key := fmt.Sprintf(public.RedisKeyTip, userId)
	rdb := connector.GetRedis().WithContext(c)
	_, err := rdb.Set(c, key, 1, time.Duration(util.SecondsUntilMidnight())*time.Second).Result()
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[TipConfirm] [Set] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, nil)
}
