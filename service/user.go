package service

import (
	"errors"
	"financia/models"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RegisterReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResp struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Register(c *gin.Context) {
	req := &RegisterReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.S().Error("[Register] [ShouldBindJSON] [err] = ", err)
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}
	// 邮箱格式验证
	if !util.CheckEmailFormat(req.Email) {
		zap.S().Error("[Register] 邮箱格式不正确")
		util.FailRespWithCode(c, util.ReqDataError)
		return
	}

	// 判断邮箱是否已经存在
	userInfoRepo := models.NewUserInfoRepo(c)
	_, err = userInfoRepo.GetUserInfoByEmail(req.Email)
	if err == nil {
		zap.S().Error("[Register] [GetUserInfoByEmail] [err] = ", err)
		util.FailRespWithCode(c, util.ReqDataError)
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.S().Error("[Register] [GetUserInfoByEmail] [err] = ", err)
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 创建用户 密码md5加密
	user := &models.UserInfo{
		Name:     req.Name,
		Email:    req.Email,
		Password: util.GetMD5(req.Password),
	}
	id, err := userInfoRepo.CreateUser(user)
	if err != nil {
		zap.S().Error("[Register] [CreateUser] [err] = ", err)
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	// 生成token
	token, err := util.GenerateJWT(id)
	if err != nil {
		zap.S().Error("[Register] [GenerateJWT] [err] = ", err)
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, &RegisterResp{
		Id:    id,
		Token: token,
		Name:  user.Name,
		Email: user.Email,
	})
	return
}

type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Login(c *gin.Context) {
	req := &LoginReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.S().Error("[Login] [ShouldBindJSON] [err] = ", err)
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		return
	}

	// 查询用户
	userInfoRepo := models.NewUserInfoRepo(c)
	user, err := userInfoRepo.GetUserInfoByEmail(req.Email)
	if err != nil {
		zap.S().Error("[Login] [GetUserInfoByEmail] [err] = ", err)
		util.FailRespWithCode(c, util.ReqDataError)
		return
	}

	// 判断密码是否正确
	if user.Password != util.GetMD5(req.Password) {
		zap.S().Error("[Login] 密码错误")
		util.FailRespWithCode(c, util.ReqDataError)
		return
	}

	// 生成token
	token, err := util.GenerateJWT(user.Id)
	if err != nil {
		zap.S().Error("[Login] [GenerateJWT] [err] = ", err)
		util.FailRespWithCode(c, util.InternalServerError)
		return
	}

	util.SuccessResp(c, &LoginResp{
		Id:    user.Id,
		Token: token,
		Name:  user.Name,
		Email: user.Email,
	})
}
