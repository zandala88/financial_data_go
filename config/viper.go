package config

import (
	"github.com/spf13/viper"
	"log"
)

var Configs Config

type Config struct {
	MySQL   MySQLConfig
	Redis   RedisConfig
	Auth    AuthConfig
	App     AppConfig
	Alpha   AlphaConfig
	Logger  LoggerConfig
	Email   EmailConfig
	TuShare TuShareConfig
	Python  PythonConfig
	Spark   SparkConfig
}

type MySQLConfig struct {
	Port                      string
	Host                      string
	Username                  string
	Password                  string
	Database                  string
	Charset                   string
	ParseTime                 string
	Loc                       string
	IgnoreRecordNotFoundError bool
	LogLevel                  int
	SlowThreshold             int
}

type AuthConfig struct {
	AccessSecret string
	AccessExpire int64
}

type RedisConfig struct {
	Addr         string
	Password     string
	Db           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
}

type AppConfig struct {
	IP   string // 应用程序 IP 地址
	Port string // HTTP 服务器端口
	Salt string // 密码加盐
}

type AlphaConfig struct {
	ApiKey string
}

type LoggerConfig struct {
	Type string
}

type EmailConfig struct {
	Server   string
	Port     int
	User     string
	Password string
}

type TuShareConfig struct {
	Token string
}

type PythonConfig struct {
	Url string
}

type SparkConfig struct {
	Password string
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig() //加载配置文件
	if err != nil {
		log.Println("viper.ReadInConfig() failed, err:", err)
		return
	}
	err = viper.Unmarshal(&Configs)
	if err != nil {
		log.Println("viper.Unmarshal() failed, err:", err)
		return
	}
}
