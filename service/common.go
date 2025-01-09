package service

import (
	"financia/models"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Check(c *gin.Context) {
	util.SuccessResp(c, nil)
}

type FeedbackReq struct {
	Type    int    `form:"type" binding:"required"`
	Content string `form:"content" binding:"required"`
}

func Feedback(c *gin.Context) {
	var req FeedbackReq
	err := c.ShouldBind(&req)
	if err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[Feedback] 参数错误", err)
		return
	}

	userId := util.GetUid(c)

	feedbackRepo := models.NewUserFeedbackRepo(c)
	err = feedbackRepo.Create(&models.UserFeedback{
		Type:    req.Type,
		Content: req.Content,
		UserId:  userId,
	})
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[Feedback] feedbackRepo.Create err = ", err)
		return
	}

	util.SuccessResp(c, nil)
}
