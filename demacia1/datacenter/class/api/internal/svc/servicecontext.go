package svc

import (
	"demacia/common/datacenter"
	"demacia/datacenter/class/api/internal/config"
	"demacia/datacenter/class/model"
	"demacia/datacenter/member/rpc/memberclient"
	"demacia/datacenter/organization/rpc/organizationclient"
	"demacia/datacenter/subject/rpc/subjectclient"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config            config.Config
	StageModel        *model.StageModel
	GradeModel        *model.GradeModel
	ClassModel        *model.ClassModel
	ClassTeacherModel *model.ClassTeacherModel
	TeachModel        *model.ClassSubjectTeacherModel
	KqPusher          *kq.Pusher
	MemberRpc         memberclient.Member
	OrgRpc            organizationclient.Organization
	SubjectRpc        subjectclient.Subject
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	kqPusher := kq.NewPusher(c.Brokers, datacenter.Class)
	return &ServiceContext{
		Config:            c,
		StageModel:        model.NewStageModel(conn, cacheRedis),
		GradeModel:        model.NewGradeModel(conn, cacheRedis),
		ClassModel:        model.NewClassModel(conn, cacheRedis),
		TeachModel:        model.NewClassSubjectTeacherModel(conn, cacheRedis),
		ClassTeacherModel: model.NewClassTeacherModel(conn, cacheRedis),
		KqPusher:          kqPusher,
		MemberRpc:         memberclient.NewMember(zrpc.MustNewClient(c.MemberRpc)),
		OrgRpc:            organizationclient.NewOrganization(zrpc.MustNewClient(c.OrgRpc)),
		SubjectRpc:        subjectclient.NewSubject(zrpc.MustNewClient(c.SubjectRpc)),
	}
}
