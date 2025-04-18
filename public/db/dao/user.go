package dao

import (
	"context"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/model"
)

// CreateUser 创建用户
func CreateUser(ctx context.Context, email, username, password string) error {
	return connector.GetDB().WithContext(ctx).Create(&model.UserInfo{
		Email:    email,
		Username: username,
		Password: public.GenerateMD5Hash(password),
	}).Error
}

func GetUser(ctx context.Context, userId int64) (*model.UserInfo, error) {
	var user *model.UserInfo
	err := connector.GetDB().WithContext(ctx).Where("f_id = ?", userId).First(&user).Error
	return user, err
}

func GetUserId(ctx context.Context, email string) int64 {
	var user *model.UserInfo
	connector.GetDB().WithContext(ctx).Where("f_email = ?", email).First(&user)
	return user.Id
}
