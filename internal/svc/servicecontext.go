// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"fmt"
	"sync"

	"mysurl1/internal/config"

	goredis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	serviceContext     *ServiceContext
	serviceContextOnce sync.Once
)

type ServiceContext struct {
	Config config.Config
	DB     sqlx.SqlConn
	Redis  *goredis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	serviceContextOnce.Do(func() {
		serviceContext = &ServiceContext{
			Config: c,
			DB:     newMySQL(c.MySQL),
			Redis:  newRedis(c.Redis),
		}
	})

	return serviceContext
}

func newMySQL(c config.MySQLConf) sqlx.SqlConn {
	if c.Host == "" || c.Port == 0 || c.User == "" || c.Database == "" {
		return nil
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)

	return sqlx.NewMysql(dsn)
}

func newRedis(c config.RedisConf) *goredis.Client {
	if c.Host == "" || c.Port == 0 {
		return nil
	}

	return goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
		DB:       c.DB,
	})
}
