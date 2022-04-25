package svc

import (
	"demacia/datacenter/class/rpc/classclient"
	"demacia/datacenter/member/rpc/memberclient"
	"demacia/datacenter/subject/model"
	"demacia/datacenter/subject/rmq/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	MemberRpc           memberclient.Member
	ClassRpc            classclient.Class
	SubjectGradeModel   *model.SubjectGradeModel
	SubjectTeacherModel *model.SubjectTeacherModel
	SubjectModel        *model.SubjectModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:              c,
		SubjectGradeModel:   model.NewSubjectGradeModel(conn, cacheRedis),
		SubjectTeacherModel: model.NewSubjectTeacherModel(conn, cacheRedis),
		SubjectModel:        model.NewSubjectModel(conn, cacheRedis),
		MemberRpc:           memberclient.NewMember(zrpc.MustNewClient(c.MemberRpc)),
		ClassRpc:            classclient.NewClass(zrpc.MustNewClient(c.ClassRpc)),
	}
}
