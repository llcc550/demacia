package svc

import (
	"demacia/datacenter/coursetable/model"
	"demacia/datacenter/coursetable/rpc/internal/config"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config                 config.Config
	CourseTableDeployModel *model.CourseTableDeployModel
	CourseTableModel       *model.CourseTableModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:                 c,
		CourseTableDeployModel: model.NewCourseTableDeployModel(conn, cacheRedis),
		CourseTableModel:       model.NewCourseTableModel(conn, cacheRedis),
	}
}
