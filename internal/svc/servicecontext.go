// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"fmt"
	"sync"

	"mysurl1/internal/config"
	"mysurl1/internal/dao"
	codestrategy "mysurl1/internal/logic/code_strategy"

	goredis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	serviceContext     *ServiceContext
	serviceContextOnce sync.Once
)

type ServiceContext struct {
	Config       config.Config
	DB           sqlx.SqlConn
	Redis        *goredis.Client
	ShortLinkDAO *dao.ShortLinkDAO
	CodeManager  *codestrategy.CodeManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	serviceContextOnce.Do(func() {
		serviceContext = &ServiceContext{
			Config: c,
			DB:     newMySQL(c.MySQL),
			Redis:  newRedis(c.Redis),
		}
		serviceContext.ShortLinkDAO = dao.NewShortLinkDAO(serviceContext.DB)
		serviceContext.CodeManager = mustNewCodeManager(c.Short, serviceContext.ShortLinkDAO)
	})

	return serviceContext
}

func mustNewCodeManager(c config.ShortConf, shortLinkDAO *dao.ShortLinkDAO) *codestrategy.CodeManager {
	manager := codestrategy.NewCodeManager(c.Provider)
	manager.Register(codestrategy.NewMySQLAutoIncrementGenerator(shortLinkDAO))
	manager.Register(codestrategy.NewRedisIncrGenerator(serviceContext.Redis))
	manager.Register(codestrategy.NewSnowflakeGenerator())

	if _, err := manager.Get(c.Provider); err != nil && c.Provider != "" {
		panic(err)
	}
	if _, err := manager.Get(""); err != nil {
		panic(err)
	}

	return manager
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
