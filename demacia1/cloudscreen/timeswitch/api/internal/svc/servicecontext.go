package svc

import (
	"demacia/cloudscreen/timeswitch/api/internal/config"
	"demacia/cloudscreen/timeswitch/api/internal/middleware"
	"demacia/cloudscreen/timeswitch/model"
	"demacia/datacenter/common/rpc/commonclient"
	"demacia/datacenter/databus/rpc/databusclient"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	Log                   rest.Middleware
	TimeSwitchConfigModel *model.TimeSwitchConfigModel
	TimeSwitchDateModel   *model.TimeSwitchDateModel
	TimeSwitchModel       *model.TimeSwitchModel
	CommonRpc             commonclient.Common
	DataBusRpc            databusclient.Databus
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	dataBusRpc := databusclient.NewDatabus(zrpc.MustNewClient(c.DataBusRpc))
	return &ServiceContext{
		Config:                c,
		Log:                   middleware.NewLogMiddleware(dataBusRpc).Handle,
		TimeSwitchConfigModel: model.NewTimeSwitchConfigModel(conn, cacheRedis),
		TimeSwitchDateModel:   model.NewTimeSwitchDateModel(conn, cacheRedis),
		TimeSwitchModel:       model.NewTimeSwitchModel(conn, cacheRedis),
		CommonRpc:             commonclient.NewCommon(zrpc.MustNewClient(c.CommonRpc)),
		DataBusRpc:            dataBusRpc,
	}
}
