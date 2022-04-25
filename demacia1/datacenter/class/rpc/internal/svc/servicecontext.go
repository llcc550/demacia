package svc

import (
	"demacia/datacenter/class/model"
	"demacia/datacenter/class/rpc/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config                   config.Config
	ClassModel               *model.ClassModel
	GradeModel               *model.GradeModel
	ClassSubjectTeacherModel *model.ClassSubjectTeacherModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:                   c,
		ClassModel:               model.NewClassModel(conn, cacheRedis),
		GradeModel:               model.NewGradeModel(conn, cacheRedis),
		ClassSubjectTeacherModel: model.NewClassSubjectTeacherModel(conn, cacheRedis),
	}
}
