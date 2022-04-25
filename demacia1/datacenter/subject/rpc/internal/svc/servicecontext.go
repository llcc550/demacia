package svc

import (
	"demacia/datacenter/subject/model"
	"demacia/datacenter/subject/rpc/internal/config"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config              config.Config
	SubjectModel        *model.SubjectModel
	SubjectGradeModel   *model.SubjectGradeModel
	SubjectTeacherModel *model.SubjectTeacherModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:              c,
		SubjectModel:        model.NewSubjectModel(conn, cacheRedis),
		SubjectGradeModel:   model.NewSubjectGradeModel(conn, cacheRedis),
		SubjectTeacherModel: model.NewSubjectTeacherModel(conn, cacheRedis),
	}
}
