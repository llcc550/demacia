package svc

import (
	"demacia/datacenter/class/rpc/classclient"
	"demacia/datacenter/databus/rpc/databusclient"
	"demacia/datacenter/device/rpc/deviceclient"
	"demacia/datacenter/position/api/internal/config"
	"demacia/datacenter/position/api/internal/middleware"
	"demacia/datacenter/position/model"
	"demacia/service/websocket/rpc/websocketclient"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	PositionModel       *model.PositionModel
	PositionDeviceModel *model.PositionDeviceModel
	ClassRpc            classclient.Class
	DeviceRpc           deviceclient.Device
	WebsocketRpc        websocketclient.Websocket
	Config              config.Config
	Log                 rest.Middleware
	DataBusRpc          databusclient.Databus
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	dataBusRpc := databusclient.NewDatabus(zrpc.MustNewClient(c.DataBusRpc))
	return &ServiceContext{
		PositionModel:       model.NewPositionModel(conn, cacheRedis),
		PositionDeviceModel: model.NewPositionDeviceModel(conn, cacheRedis),
		ClassRpc:            classclient.NewClass(zrpc.MustNewClient(c.ClassRpc)),
		DeviceRpc:           deviceclient.NewDevice(zrpc.MustNewClient(c.DeviceRpc)),
		WebsocketRpc:        websocketclient.NewWebsocket(zrpc.MustNewClient(c.WebsocketRpc)),
		Config:              c,
		DataBusRpc:          dataBusRpc,
		Log:                 middleware.NewLogMiddleware(dataBusRpc).Handle,
	}
}
