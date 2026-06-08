// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Stat       StatConf       `json:",optional"`
	MySQL      MySQLConf      `json:",optional"`
	Redis      RedisConf      `json:",optional"`
	Auth       AuthConf       `json:",optional"`
	Short      ShortConf      `json:",optional"`
	VisitFlush VisitFlushConf `json:",optional"`
}

type StatConf struct {
	DisableLog     bool `json:",default=true"`
	DisableSampler bool `json:",optional"`
}

type MySQLConf struct {
	Host     string `json:",optional"`
	Port     int    `json:",optional"`
	User     string `json:",optional"`
	Password string `json:",optional"`
	Database string `json:",optional"`
}

type RedisConf struct {
	Host     string `json:",optional"`
	Port     int    `json:",optional"`
	Password string `json:",optional"`
	DB       int    `json:",optional"`
}

type AuthConf struct {
	JWTSecret      string `json:",optional"`
	ExpireSeconds  int64  `json:",default=86400"`
	PasswordPepper string `json:",optional"`
}

type ShortConf struct {
	BaseURL   string        `json:",optional"`
	Provider  string        `json:",default=mysql_auto_increment,options=mysql_auto_increment|redis_incr|snowflake"`
	Snowflake SnowflakeConf `json:",optional"`
}

type SnowflakeConf struct {
	WorkerID int64 `json:",optional"`
}

type VisitFlushConf struct {
	Interval time.Duration `json:",default=5s"`
	Batch    int64         `json:",default=100"`
}
