package config

import (
	"github.com/spf13/viper"
)

var Configs Config

type Config struct {
	MySQL    MySQLConfig
	Redis    RedisConfig
	Auth     AuthConfig
	App      AppConfig
	ETCD     ETCDConfig
	RabbitMQ RabbitMQConfig
	Logger   LoggerConfig
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
	Salt              string // 密码加盐
	IP                string // 应用程序 IP 地址
	HTTPServerPort    string // HTTP 服务器端口
	WebsocketPort     string // WebSocket 服务器端口
	RPCPort           string // RPC 服务器端口
	WorkerPoolSize    uint32 // 业务 worker 队列数量
	MaxWorkerTask     int    // 业务 worker 对应负责的任务队列最大任务存储数量
	HeartbeatTimeout  int    // 心跳超时时间（秒）
	HeartbeatInterval int    // 超时连接检测间隔（秒）
}

type ETCDConfig struct {
	Endpoints  []string // endpoints 列表
	Timeout    int      // 超时时间（秒）
	ServerList string   // 服务列表
}

type RabbitMQConfig struct {
	URL string // rabbitmq url
}

type LoggerConfig struct {
	Type string
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig() //加载配置文件
	if err != nil {
		return
	}
	err = viper.Unmarshal(&Configs)
	if err != nil {
		return
	}
}
