package svc

import (
	"demacia/common/datacenter"
	"demacia/datacenter/card/rpc/cardclient"
	"demacia/datacenter/class/rpc/classclient"
	"demacia/datacenter/student/api/internal/config"
	"demacia/datacenter/student/model"
	"demacia/service/websocket/rpc/websocketclient"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	ClassRpc     classclient.Class
	WebsocketRpc websocketclient.Websocket
	CardRpc      cardclient.Card
	StudentModel *model.StudentModel
	KqPusher     *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	kqPusher := kq.NewPusher(c.Brokers, datacenter.Student)
	return &ServiceContext{
		Config:       c,
		ClassRpc:     classclient.NewClass(zrpc.MustNewClient(c.ClassRpc)),
		WebsocketRpc: websocketclient.NewWebsocket(zrpc.MustNewClient(c.WebsocketRpc)),
		CardRpc:      cardclient.NewCard(zrpc.MustNewClient(c.CardRpc)),
		StudentModel: model.NewStudentModel(conn, cacheRedis),
		KqPusher:     kqPusher,
	}
}
