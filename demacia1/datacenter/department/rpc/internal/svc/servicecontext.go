package svc

import (
	"demacia/datacenter/department/model"
	"demacia/datacenter/department/rpc/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config                config.Config
	DepartmentModel       *model.DepartmentModel
	DepartmentMemberModel *model.DepartmentMemberModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:                c,
		DepartmentModel:       model.NewDepartmentModel(conn, cacheRedis),
		DepartmentMemberModel: model.NewDepartmentMemberModel(conn, cacheRedis),
	}
}
