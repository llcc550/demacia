package svc

import (
	"demacia/datacenter/class/rpc/classclient"
	"demacia/datacenter/coursetable/api/internal/config"
	"demacia/datacenter/coursetable/api/internal/middleware"
	"demacia/datacenter/coursetable/model"
	"demacia/datacenter/databus/rpc/databusclient"
	"demacia/datacenter/member/rpc/memberclient"
	"demacia/datacenter/organization/rpc/organizationclient"
	"demacia/datacenter/position/rpc/positionclient"
	"demacia/datacenter/student/rpc/studentclient"
	"demacia/datacenter/subject/rpc/subjectclient"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config                 config.Config
	Log                    rest.Middleware
	CourseTableModel       *model.CourseTableModel
	CourseTableDeployModel *model.CourseTableDeployModel
	TeachModel             *model.ClassSubjectTeacherModel
	MemberRpc              memberclient.Member
	PositionRpc            positionclient.Position
	StudentRpc             studentclient.Student
	DataBusRpc             databusclient.Databus
	SubjectRpc             subjectclient.Subject
	ClassRpc               classclient.Class
	OrgRpc                 organizationclient.Organization
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	dataBusRpc := databusclient.NewDatabus(zrpc.MustNewClient(c.DataBusRpc))
	return &ServiceContext{
		Config:                 c,
		TeachModel:             model.NewClassSubjectTeacherModel(conn, cacheRedis),
		CourseTableModel:       model.NewCourseTableModel(conn, cacheRedis),
		CourseTableDeployModel: model.NewCourseTableDeployModel(conn, cacheRedis),
		PositionRpc:            positionclient.NewPosition(zrpc.MustNewClient(c.PositionRpc)),
		MemberRpc:              memberclient.NewMember(zrpc.MustNewClient(c.MemberRpc)),
		StudentRpc:             studentclient.NewStudent(zrpc.MustNewClient(c.StudentRpc)),
		OrgRpc:                 organizationclient.NewOrganization(zrpc.MustNewClient(c.OrgRpc)),
		SubjectRpc:             subjectclient.NewSubject(zrpc.MustNewClient(c.SubjectRpc)),
		ClassRpc:               classclient.NewClass(zrpc.MustNewClient(c.ClassRpc)),
		DataBusRpc:             dataBusRpc,
		Log:                    middleware.NewLogMiddleware(dataBusRpc).Handle,
	}
}
