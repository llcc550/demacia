package svc

import (
	"demacia/datacenter/member/model"
	"demacia/datacenter/member/rmq/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config      config.Config
	Redis       *redis.Redis
	MemberModel *model.MemberModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:      c,
		Redis:       cacheRedis,
		MemberModel: model.NewMemberModel(conn, cacheRedis),
	}
}
