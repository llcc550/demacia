package svc

import (
	"demacia/datacenter/position/model"
	"demacia/datacenter/position/rmq/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config              config.Config
	Redis               *redis.Redis
	PositionModel       *model.PositionModel
	PositionDeviceModel *model.PositionDeviceModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:              c,
		Redis:               cacheRedis,
		PositionModel:       model.NewPositionModel(conn, cacheRedis),
		PositionDeviceModel: model.NewPositionDeviceModel(conn, cacheRedis),
	}
}
