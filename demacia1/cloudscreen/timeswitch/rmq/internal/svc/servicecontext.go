package svc

import (
	"demacia/cloudscreen/timeswitch/model"
	"demacia/cloudscreen/timeswitch/rmq/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config                config.Config
	Redis                 *redis.Redis
	TimeSwitchModel       *model.TimeSwitchModel
	TimeSwitchConfigModel *model.TimeSwitchConfigModel
	TimeSwitchDateModel   *model.TimeSwitchDateModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:                c,
		Redis:                 cacheRedis,
		TimeSwitchModel:       model.NewTimeSwitchModel(conn, cacheRedis),
		TimeSwitchDateModel:   model.NewTimeSwitchDateModel(conn, cacheRedis),
		TimeSwitchConfigModel: model.NewTimeSwitchConfigModel(conn, cacheRedis),
	}
}
