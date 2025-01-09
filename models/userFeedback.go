package models

import (
	"context"
	"financia/public"
	"gorm.io/gorm"
	"time"
)

type UserFeedback struct {
	Id      int       `gorm:"column:f_id;primaryKey;autoIncrement" json:"id"`               // 主键，自增 ID
	UserId  int64     `gorm:"column:f_user_id;not null" json:"userId"`                      // 关联用户表的用户 ID
	Content string    `gorm:"column:f_content;type:text;not null" json:"content"`           // 用户的反馈内容
	Time    time.Time `gorm:"column:f_time;not null;default:CURRENT_TIMESTAMP" json:"time"` // 反馈时间
	Status  int       `gorm:"column:f_status;not null;default:0" json:"status"`             // 反馈状态：0 未处理，1 已处理
	Type    int       `gorm:"column:f_type;not null" json:"type"`                           // 反馈类型：1 Bug，2 反馈，3 功能
}

// TableName 设置表名
func (UserFeedback) TableName() string {
	return "t_user_feedback"
}

type UserFeedbackRepo struct {
	db      *gorm.DB
	context context.Context
}

func NewUserFeedbackRepo(ctx context.Context) *UserFeedbackRepo {
	return &UserFeedbackRepo{
		db:      public.DB.WithContext(ctx),
		context: ctx,
	}
}

func (r *UserFeedbackRepo) Create(feedback *UserFeedback) error {
	err := r.db.Create(feedback).Error
	return err
}
