package connector

import (
	"financia/config"
	"financia/public/db/model"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/sharding"
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

	// 配置连接池
	sqlDB, err := mysql.DB()
	if err != nil {
		zap.S().Error("[init] [Get DB instance] [err] = ", err.Error())
		panic(err)
	}
	sqlDB.SetMaxOpenConns(100)                 // 设置最大连接数
	sqlDB.SetMaxIdleConns(50)                  // 设置最大空闲连接数
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // 设置连接最大生命周期

	// 注册分表插件
	err = mysql.Use(sharding.Register(sharding.Config{
		ShardingKey:         "f_ts_code",          // 分片键
		NumberOfShards:      20,                   // 分片数量
		PrimaryKeyGenerator: sharding.PKSnowflake, // 使用 Snowflake 算法生成主键
	}, model.StockData{}, model.FundData{})) // 注册需要分表的表
	if err != nil {
		panic(fmt.Sprintf("failed to register sharding plugin: %v", err))
	}

	db = mysql
}

func GetDB() *gorm.DB {
	return db
}
