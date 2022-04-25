package svc

import (
	"demacia/datacenter/card/rpc/cardclient"
	"demacia/datacenter/databus/rpc/databusclient"
	"demacia/datacenter/department/rpc/departmentclient"
	"demacia/datacenter/member/api/internal/config"
	"demacia/datacenter/member/api/internal/middleware"
	"demacia/datacenter/member/model"
	"demacia/datacenter/organization/rpc/organizationclient"
	"demacia/service/websocket/rpc/websocketclient"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	Log             rest.Middleware
	MemberModel     *model.MemberModel
	OrganizationRpc organizationclient.Organization
	WebsocketRpc    websocketclient.Websocket
	CardRpc         cardclient.Card
	DepartmentRpc   departmentclient.Department
	DataBusRpc      databusclient.Databus
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	dataBusRpc := databusclient.NewDatabus(zrpc.MustNewClient(c.DataBusRpc))
	return &ServiceContext{
		Config:          c,
		Log:             middleware.NewLogMiddleware(dataBusRpc).Handle,
		MemberModel:     model.NewMemberModel(conn, cacheRedis),
		WebsocketRpc:    websocketclient.NewWebsocket(zrpc.MustNewClient(c.WebsocketRpc)),
		OrganizationRpc: organizationclient.NewOrganization(zrpc.MustNewClient(c.OrganizationRpc)),
		CardRpc:         cardclient.NewCard(zrpc.MustNewClient(c.CardRpc)),
		DepartmentRpc:   departmentclient.NewDepartment(zrpc.MustNewClient(c.DepartmentRpc)),
		DataBusRpc:      dataBusRpc,
	}
}
