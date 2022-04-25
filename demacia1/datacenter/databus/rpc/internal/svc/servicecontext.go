package svc

import (
	"demacia/common/datacenter"
	"demacia/datacenter/databus/model"
	"demacia/datacenter/databus/rpc/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
)

type ServiceContext struct {
	Config   config.Config
	KqPusher *kq.Pusher
	LogModel model.LogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	kqPusher := kq.NewPusher(c.Brokers, datacenter.Kafka)
	return &ServiceContext{
		Config:   c,
		KqPusher: kqPusher,
		LogModel: model.NewLogModel(c.Mongo.Url, c.Mongo.Collection.Log),
	}
}
