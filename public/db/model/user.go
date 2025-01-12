package model

import "time"

type UserInfo struct {
	Id        int64     `gorm:"column:f_id;primaryKey;autoIncrement;comment:用户ID"`
	Email     string    `gorm:"column:f_email;size:50;not null;unique;comment:邮箱"`
	Username  string    `gorm:"column:f_username;size:50;not null;unique;comment:用户名"`
	Password  string    `gorm:"column:f_password;size:255;not null;comment:密码"`
	CreatedAt time.Time `gorm:"column:f_created_at;autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time `gorm:"column:f_updated_at;autoUpdateTime;comment:更新时间"`
	Status    int       `gorm:"column:f_status;default:0;comment:用户状态: 0=正常, 1=禁用"`
}

// TableName 设置表名
func (UserInfo) TableName() string {
	return "t_user_info"
}
