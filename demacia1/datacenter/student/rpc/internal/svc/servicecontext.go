package svc

import (
	"demacia/datacenter/student/model"
	"demacia/datacenter/student/rpc/internal/config"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config       config.Config
	StudentModel *model.StudentModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:       c,
		StudentModel: model.NewStudentModel(conn, cacheRedis),
	}
}
