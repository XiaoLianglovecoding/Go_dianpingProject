package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	// Server 保存 Web 服务自己的配置，比如监听端口。
	Server ServerConfig `mapstructure:"server"`
	// MySQL 保存连接数据库需要的配置。
	MySQL MySQLConfig `mapstructure:"mysql"`
	// Redis 保存连接 Redis 需要的配置。
	Redis RedisConfig `mapstructure:"redis"`
}

type ServerConfig struct {
	// Port 对应 configs/config.yaml 里的 server.port。
	Port int `mapstructure:"port"`
}

type MySQLConfig struct {
	// Host 是数据库地址，本机一般是 127.0.0.1。
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"` // Port 是数据库端口，本机 MySQL 常见是 3306，Docker 映射可能是 3307。
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	// Host 是 Redis 地址；Database 是 Redis 逻辑库编号，Java 版项目使用的是 1。
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

// Load 读取配置文件，并把 yaml 里的配置反序列化成 Config 结构体。
//
// 例如 configs/config.yaml 中的:
//
//	server:
//	  port: 8081
//
// 会被读取到 cfg.Server.Port。
func Load(path string) (*Config, error) {
	// viper 是 Go 里常用的配置读取库；这里新建一个独立实例，避免全局配置互相影响。
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	// 允许用环境变量覆盖配置，例如 SERVER_PORT 可以覆盖 server.port。
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 如果配置文件没写 server.port，就默认使用 8081，和 nginx 代理配置保持一致。
	v.SetDefault("server.port", 8081)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
