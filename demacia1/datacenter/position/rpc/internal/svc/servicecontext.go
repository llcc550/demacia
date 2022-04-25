package svc

import (
	"demacia/datacenter/position/model"
	"demacia/datacenter/position/rpc/internal/config"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config        config.Config
	PositionModel *model.PositionModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:        c,
		PositionModel: model.NewPositionModel(conn, cacheRedis),
	}
}
