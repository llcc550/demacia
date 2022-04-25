package svc

import (
	"demacia/datacenter/department/model"
	"demacia/datacenter/department/rmq/internal/config"
	"demacia/datacenter/member/rpc/memberclient"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	DepartmentModel       *model.DepartmentModel
	DepartmentMemberModel *model.DepartmentMemberModel
	MemberRpc             memberclient.Member
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:                c,
		DepartmentModel:       model.NewDepartmentModel(conn, cacheRedis),
		DepartmentMemberModel: model.NewDepartmentMemberModel(conn, cacheRedis),
		MemberRpc:             memberclient.NewMember(zrpc.MustNewClient(c.MemberRpc)),
	}
}
