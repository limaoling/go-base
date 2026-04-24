package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	DBName          string
	Charset         string
	ParseTime       bool
	Loc             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

type JWTConfig struct {
	Secret     string
	ExpireHour int
}

var GlobalConfig *Config

func InitConfig(configPath string) error {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	GlobalConfig = &Config{
		Server: ServerConfig{
			Port: v.GetInt("server.port"),
		},
		Database: DatabaseConfig{
			Host:            v.GetString("database.host"),
			Port:            v.GetInt("database.port"),
			Username:        v.GetString("database.username"),
			Password:        v.GetString("database.password"),
			DBName:          v.GetString("database.dbname"),
			Charset:         v.GetString("database.charset"),
			ParseTime:       v.GetBool("database.parse_time"),
			Loc:             v.GetString("database.loc"),
			MaxIdleConns:    v.GetInt("database.max_idle_conns"),
			MaxOpenConns:    v.GetInt("database.max_open_conns"),
			ConnMaxLifetime: v.GetInt("database.conn_max_lifetime"),
		},
		JWT: JWTConfig{
			Secret:     v.GetString("jwt.secret"),
			ExpireHour: v.GetInt("jwt.expire_hour"),
		},
	}

	return nil
}

// GetDSN 返回数据库 DSN 连接字符串
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.DBName,
		d.Charset,
		d.ParseTime,
		d.Loc,
	)
}

// GetDSNWithoutDB 返回不带数据库名的 DSN 连接字符串（用于创建数据库）
func (d *DatabaseConfig) GetDSNWithoutDB() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=%t&loc=%s",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Charset,
		d.ParseTime,
		d.Loc,
	)
}
