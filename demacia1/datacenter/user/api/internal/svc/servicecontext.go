package svc

import (
	"demacia/datacenter/user/api/internal/config"
	"demacia/datacenter/user/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config    config.Config
	UserModel *model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(conn, cacheRedis),
	}
}
