package connector

import (
	"financia/config"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var db *gorm.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		config.Configs.MySQL.Username, config.Configs.MySQL.Password, config.Configs.MySQL.Host,
		config.Configs.MySQL.Port, config.Configs.MySQL.Database, config.Configs.MySQL.Charset,
		config.Configs.MySQL.ParseTime, config.Configs.MySQL.Loc)

	zap.S().Debug("[init] [mysql] [dsn] = ", dsn)

	mysql, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		QueryFields:            true, //打印sql
		Logger: &CustomLogger{
			logLevel:                  logger.LogLevel(config.Configs.MySQL.LogLevel),                       // 日志等级
			ignoreRecordNotFoundError: config.Configs.MySQL.IgnoreRecordNotFoundError,                       // true 忽略 ErrRecordNotFound 错误
			slowThreshold:             time.Duration(config.Configs.MySQL.SlowThreshold) * time.Millisecond, // 慢查询阈值
		},
	})
	if err != nil {
		zap.S().Error("[init] [gorm.Open] [err] = ", err.Error())
		panic(err)
	}

	db = mysql
}

func GetDB() *gorm.DB {
	return db
}
