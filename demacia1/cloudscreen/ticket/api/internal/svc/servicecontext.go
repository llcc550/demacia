package svc

import (
	"demacia/cloudscreen/ticket/api/internal/config"
	"demacia/cloudscreen/ticket/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config      config.Config
	TicketModel *model.TicketModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:      c,
		TicketModel: model.NewTicketModel(conn, cacheRedis),
	}
}
