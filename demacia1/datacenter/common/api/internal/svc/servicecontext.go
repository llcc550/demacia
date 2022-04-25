package svc

import (
	"demacia/datacenter/common/api/internal/config"
	"demacia/datacenter/common/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config      config.Config
	AreaModel   *model.AreaModel
	EthnicModel *model.EthnicModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:      c,
		AreaModel:   model.NewAreaModel(conn, cacheRedis),
		EthnicModel: model.NewEthnicModel(conn, cacheRedis),
	}
}
