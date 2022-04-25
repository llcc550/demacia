package svc

import (
	"demacia/common/datacenter"
	"demacia/datacenter/class/rpc/classclient"
	"demacia/datacenter/parent/api/internal/config"
	"demacia/datacenter/parent/model"
	"demacia/datacenter/student/rpc/studentclient"
	"demacia/service/websocket/rpc/websocketclient"
	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	ClassRpc           classclient.Class
	StudentRpc         studentclient.Student
	ParentModel        *model.ParentModel
	StudentParentModel *model.StudentParentModel
	WebsocketRpc       websocketclient.Websocket
	KqPusher           *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	kqPusher := kq.NewPusher(c.Brokers, datacenter.Parent)
	return &ServiceContext{
		Config:             c,
		ClassRpc:           classclient.NewClass(zrpc.MustNewClient(c.ClassRpc)),
		StudentRpc:         studentclient.NewStudent(zrpc.MustNewClient(c.StudentRpc)),
		ParentModel:        model.NewParentModel(conn, cacheRedis),
		StudentParentModel: model.NewStudentParentModel(conn, cacheRedis),
		WebsocketRpc:       websocketclient.NewWebsocket(zrpc.MustNewClient(c.WebsocketRpc)),
		KqPusher:           kqPusher,
	}
}
