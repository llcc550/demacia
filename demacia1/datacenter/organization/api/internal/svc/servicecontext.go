package svc

import (
	"demacia/datacenter/common/rpc/commonclient"
	"demacia/datacenter/databus/rpc/databusclient"
	"demacia/datacenter/member/rpc/memberclient"
	"demacia/datacenter/organization/api/internal/config"
	"demacia/datacenter/organization/api/internal/middleware"
	"demacia/datacenter/organization/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config            config.Config
	Log               rest.Middleware
	OrganizationModel *model.OrganizationModel
	CommonRpc         commonclient.Common
	MemberRpc         memberclient.Member
	DataBusRpc        databusclient.Databus
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	dataBusRpc := databusclient.NewDatabus(zrpc.MustNewClient(c.DataBusRpc))
	return &ServiceContext{
		Config:            c,
		Log:               middleware.NewLogMiddleware(dataBusRpc).Handle,
		OrganizationModel: model.NewOrganizationModel(conn, cacheRedis),
		CommonRpc:         commonclient.NewCommon(zrpc.MustNewClient(c.CommonRpc)),
		MemberRpc:         memberclient.NewMember(zrpc.MustNewClient(c.MemberRpc)),
		DataBusRpc:        dataBusRpc,
	}
}
