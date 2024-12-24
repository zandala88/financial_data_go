package models

import (
	"context"
	"financia/public"
	"gorm.io/gorm"
	"time"
)

type UserInfo struct {
	Id            int64      `gorm:"primaryKey;column:f_id" json:"id"`                                                 // 自增id，数据库列名为f_id，JSON字段名为id
	Name          string     `gorm:"default:'';column:f_name" json:"name"`                                             // 用户名，默认值为空字符串，数据库列名为f_name，JSON字段名为name
	Email         string     `gorm:"default:'';column:f_email" json:"email"`                                           // 邮箱，默认值为空字符串，数据库列名为f_email，JSON字段名为email
	Password      string     `gorm:"not null;column:f_password" json:"password"`                                       // 密码，不能为空，数据库列名为f_password，JSON字段名为password
	CreateTime    time.Time  `gorm:"default:current_timestamp;column:f_create_time" json:"create_time"`                // 创建时间，数据库列名为f_create_time，JSON字段名为create_time
	UpdateTime    time.Time  `gorm:"default:current_timestamp;autoUpdateTime;column:f_update_time" json:"update_time"` // 更新时间，数据库列名为f_update_time，JSON字段名为update_time
	LastLoginTime *time.Time `gorm:"default:null;column:f_last_login_time" json:"last_login_time"`                     // 最后登录时间，数据库列名为f_last_login_time，JSON字段名为last_login_time
}

// TableName 指定表名为 t_user_info
func (UserInfo) TableName() string {
	return "t_user_info"
}

type UserInfoRepo struct {
	db  *gorm.DB
	ctx context.Context
}

func NewUserInfoRepo(ctx context.Context) *UserInfoRepo {
	return &UserInfoRepo{
		db:  public.DB.WithContext(ctx),
		ctx: ctx,
	}
}

func (r *UserInfoRepo) CreateUser(user *UserInfo) (int64, error) {
	err := r.db.Create(user).Error
	return user.Id, err
}

func (r *UserInfoRepo) GetUserInfoByEmail(email string) (*UserInfo, error) {
	var userInfo UserInfo
	err := r.db.Where("f_email = ?", email).First(&userInfo).Error
	return &userInfo, err
}
