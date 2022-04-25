package svc

import (
	"demacia/datacenter/class/model"
	"demacia/datacenter/class/rmq/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config     config.Config
	Redis      *redis.Redis
	ClassModel *model.ClassModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:     c,
		Redis:      cacheRedis,
		ClassModel: model.NewClassModel(conn, cacheRedis),
	}
}
