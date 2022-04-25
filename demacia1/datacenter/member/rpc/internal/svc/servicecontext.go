package svc

import (
	"demacia/datacenter/databus/rpc/databusclient"
	"demacia/datacenter/member/model"
	"demacia/datacenter/member/rpc/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	MemberModel *model.MemberModel
	DataBusRpc  databusclient.Databus
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:      c,
		MemberModel: model.NewMemberModel(conn, cacheRedis),
		DataBusRpc:  databusclient.NewDatabus(zrpc.MustNewClient(c.DataBusRpc)),
	}
}
