package svc

import (
	"demacia/datacenter/device/rpc/deviceclient"
	"demacia/datacenter/photowall/api/internal/config"
	"demacia/datacenter/photowall/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config           config.Config
	PhotoFolderModel *model.PhotoFolderModel
	PhotosModel      *model.PhotosModel
	DevicePhotoModel *model.DevicePhotoModel
	DeviceRpc        deviceclient.Device
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:           c,
		PhotoFolderModel: model.NewPhotoFolderModel(conn, cacheRedis),
		PhotosModel:      model.NewPhotosModel(conn, cacheRedis),
		DevicePhotoModel: model.NewDevicePhotoModel(conn, cacheRedis),
		DeviceRpc:        deviceclient.NewDevice(zrpc.MustNewClient(c.DeviceRpc)),
	}
}
