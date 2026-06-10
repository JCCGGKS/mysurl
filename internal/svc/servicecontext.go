// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"fmt"
	"sync"

	"mysurl1/internal/config"
	"mysurl1/internal/dao"
	codestrategy "mysurl1/internal/logic/code_strategy"
	"mysurl1/internal/utils"

	goredis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
)

var (
	serviceContext     *ServiceContext
	serviceContextOnce sync.Once
)

type ServiceContext struct {
	Config              config.Config
	DB                  sqlx.SqlConn
	Redis               *goredis.Client
	ShortLinkCache      *dao.ShortLinkCache
	ShortLinkDAO        *dao.ShortLinkDAO
	UserDAO             *dao.UserDAO
	UserOperationLogDAO *dao.UserOperationLogDAO
	CodeManager         *codestrategy.CodeManager
	FlightGroup         syncx.SingleFlight
}

func NewServiceContext(c config.Config) *ServiceContext {
	serviceContextOnce.Do(func() {
		serviceContext = &ServiceContext{
			Config:      c,
			DB:          newMySQL(c.MySQL),
			Redis:       newRedis(c.Redis),
			FlightGroup: syncx.NewSingleFlight(),
		}
		serviceContext.ShortLinkCache = dao.NewShortLinkCache(serviceContext.Redis)
		serviceContext.ShortLinkDAO = dao.NewShortLinkDAO(serviceContext.DB)
		serviceContext.UserDAO = dao.NewUserDAO(serviceContext.DB)
		serviceContext.UserOperationLogDAO = dao.NewUserOperationLogDAO(serviceContext.DB)
		serviceContext.CodeManager = mustNewCodeManager(c.Short, serviceContext.ShortLinkDAO)
		utils.StartVisitFlushWorker(serviceContext.DB, serviceContext.ShortLinkCache, c.VisitFlush)
	})

	return serviceContext
}

func mustNewCodeManager(short config.ShortConf, shortLinkDAO *dao.ShortLinkDAO) *codestrategy.CodeManager {
	manager := codestrategy.NewCodeManager(short.Provider)
	manager.Register(codestrategy.NewRedisIncrGenerator(serviceContext.Redis))
	snowflakeGenerator, err := mustNewSnowflakeGenerator(short)
	if err != nil {
		panic(err)
	}
	manager.Register(snowflakeGenerator)

	if _, err := manager.Get(short.Provider); err != nil && short.Provider != "" {
		panic(err)
	}
	if _, err := manager.Get(""); err != nil {
		panic(err)
	}

	return manager
}

func mustNewSnowflakeGenerator(short config.ShortConf) (*codestrategy.SnowflakeGenerator, error) {
	if short.Provider == codestrategy.ProviderSnowflake && short.Snowflake.WorkerID == 0 {
		return nil, fmt.Errorf("short.snowflake.workerid is required when provider is snowflake")
	}

	workerID := short.Snowflake.WorkerID
	if workerID == 0 {
		workerID = 1
	}

	return codestrategy.NewSnowflakeGenerator(workerID)
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
