package svc

import (
	"demacia/cloudscreen/courserecord/api/internal/config"
	"demacia/cloudscreen/courserecord/api/internal/middleware"
	model2 "demacia/cloudscreen/courserecord/model"
	"demacia/common/datacenter"
	"demacia/datacenter/class/rpc/classclient"
	"demacia/datacenter/coursetable/rpc/coursetableclient"
	"demacia/datacenter/databus/rpc/databusclient"
	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config                  config.Config
	CourseRecordModel       *model2.CourseRecordModel
	CourseRecordConfigModel *model2.CourseRecordConfigModel
	CourseRecordCountModel  *model2.CourseRecordCountModel
	CourseRecordDateModel   *model2.CourseRecordDateModel
	CourseTableRpc          coursetableclient.Coursetable
	ClassRpc                classclient.Class
	DataBusRpc              databusclient.Databus
	KqPusher                *kq.Pusher
	Log                     rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	kqPusher := kq.NewPusher(c.Brokers, datacenter.Class)
	dataBusRpc := databusclient.NewDatabus(zrpc.MustNewClient(c.DataBusRpc))
	return &ServiceContext{
		Config:                  c,
		CourseRecordModel:       model2.NewCourseRecordModel(conn, cacheRedis),
		CourseRecordConfigModel: model2.NewCourseRecordConfigModel(conn, cacheRedis),
		CourseRecordCountModel:  model2.NewCourseRecordCountModel(conn, cacheRedis),
		CourseRecordDateModel:   model2.NewCourseRecordDateModel(conn, cacheRedis),
		CourseTableRpc:          coursetableclient.NewCoursetable(zrpc.MustNewClient(c.CourseTableRpc)),
		DataBusRpc:              dataBusRpc,
		ClassRpc:                classclient.NewClass(zrpc.MustNewClient(c.ClassRpc)),
		KqPusher:                kqPusher,
		Log:                     middleware.NewLogMiddleware(dataBusRpc).Handle,
	}
}
