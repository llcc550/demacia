package svc

import (
	"demacia/common/datacenter"
	"demacia/datacenter/device/api/internal/config"
	"demacia/datacenter/device/model"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config           config.Config
	DeviceModel      *model.DeviceModel
	GroupModel       *model.GroupModel
	DeviceGroupModel *model.DeviceGroupModel
	KqPusher         *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	kqPusher := kq.NewPusher(c.Brokers, datacenter.Class)
	return &ServiceContext{
		Config:           c,
		DeviceModel:      model.NewDeviceModel(conn, cacheRedis),
		GroupModel:       model.NewGroupModel(conn, cacheRedis),
		DeviceGroupModel: model.NewDeviceGroupModel(conn, cacheRedis),
		KqPusher:         kqPusher,
	}
}
