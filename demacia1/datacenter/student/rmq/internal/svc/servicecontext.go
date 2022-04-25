package svc

import (
	"demacia/datacenter/student/model"
	"demacia/datacenter/student/rmq/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config       config.Config
	Redis        *redis.Redis
	StudentModel *model.StudentModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:       c,
		Redis:        cacheRedis,
		StudentModel: model.NewStudentModel(conn, cacheRedis),
	}
}
