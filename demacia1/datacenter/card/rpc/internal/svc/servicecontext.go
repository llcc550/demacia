package svc

import (
	"demacia/datacenter/card/model"
	"demacia/datacenter/card/rpc/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config    config.Config
	CardModel *model.CardModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:    c,
		CardModel: model.NewCardModel(conn, cacheRedis),
	}
}
