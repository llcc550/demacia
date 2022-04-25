package svc

import (
	"demacia/datacenter/coursetable/model"
	"demacia/datacenter/coursetable/rmq/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config                 config.Config
	Redis                  *redis.Redis
	CourseTableModel       *model.CourseTableModel
	CourseTableDeployModel *model.CourseTableDeployModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:                 c,
		Redis:                  cacheRedis,
		CourseTableModel:       model.NewCourseTableModel(conn, cacheRedis),
		CourseTableDeployModel: model.NewCourseTableDeployModel(conn, cacheRedis),
	}
}
