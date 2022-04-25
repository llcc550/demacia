package svc

import (
	"demacia/datacenter/organization/model"
	"demacia/datacenter/organization/rpc/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config            config.Config
	OrganizationModel *model.OrganizationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:            c,
		OrganizationModel: model.NewOrganizationModel(conn, cacheRedis),
	}
}
