package svc

import (
	"demacia/datacenter/common/model"
	"demacia/datacenter/common/rpc/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config       config.Config
	AreaModel    *model.AreaModel
	HolidayModel *model.HolidayModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:       c,
		AreaModel:    model.NewAreaModel(conn, cacheRedis),
		HolidayModel: model.NewHolidayModel(conn, cacheRedis),
	}
}
