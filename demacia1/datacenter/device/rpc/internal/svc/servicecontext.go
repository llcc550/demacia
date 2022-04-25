package svc

import (
	"demacia/datacenter/device/model"
	"demacia/datacenter/device/rpc/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config      config.Config
	DeviceModel *model.DeviceModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:      c,
		DeviceModel: model.NewDeviceModel(conn, cacheRedis),
	}
}
