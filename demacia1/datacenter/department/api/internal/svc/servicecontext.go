package svc

import (
	"demacia/common/datacenter"
	"demacia/datacenter/department/api/internal/config"
	"demacia/datacenter/department/model"
	"demacia/datacenter/member/rpc/memberclient"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	DepartmentModel       *model.DepartmentModel
	DepartmentMemberModel *model.DepartmentMemberModel
	KqPusher              *kq.Pusher
	MemberRpc             memberclient.Member
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	kqPusher := kq.NewPusher(c.Brokers, datacenter.Department)
	return &ServiceContext{
		Config:                c,
		DepartmentModel:       model.NewDepartmentModel(conn, cacheRedis),
		DepartmentMemberModel: model.NewDepartmentMemberModel(conn, cacheRedis),
		KqPusher:              kqPusher,
		MemberRpc:             memberclient.NewMember(zrpc.MustNewClient(c.MemberRpc)),
	}
}
