package svc

import (
	"demacia/datacenter/class/rpc/classclient"
	"demacia/datacenter/databus/rpc/databusclient"
	"demacia/datacenter/member/rpc/memberclient"
	"demacia/datacenter/subject/api/internal/config"
	"demacia/datacenter/subject/api/internal/middleware"
	"demacia/datacenter/subject/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	SubjectGradeModel   *model.SubjectGradeModel
	SubjectTeacherModel *model.SubjectTeacherModel
	SubjectModel        *model.SubjectModel
	ClassRpc            classclient.Class
	MemberRpc           memberclient.Member
	DataBusRpc          databusclient.Databus
	Log                 rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	dataBusRpc := databusclient.NewDatabus(zrpc.MustNewClient(c.DataBusRpc))
	return &ServiceContext{
		Config:              c,
		Log:                 middleware.NewLogMiddleware(dataBusRpc).Handle,
		SubjectGradeModel:   model.NewSubjectGradeModel(conn, cacheRedis),
		SubjectTeacherModel: model.NewSubjectTeacherModel(conn, cacheRedis),
		SubjectModel:        model.NewSubjectModel(conn, cacheRedis),
		ClassRpc:            classclient.NewClass(zrpc.MustNewClient(c.ClassRpc)),
		MemberRpc:           memberclient.NewMember(zrpc.MustNewClient(c.MemberRpc)),
		DataBusRpc:          dataBusRpc,
	}
}
