package svc

import (
	"demacia/service/sms/model"
	"demacia/service/sms/rmq/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/limit"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config       config.Config
	HuaweiModel  *model.HuaweiModel
	TrafficModel *model.SmsTrafficModel
	Push         bool
	Limiter      *limit.PeriodLimit
}

func NewServiceContext(c config.Config) *ServiceContext {
	cacheRedis := c.CacheRedis.NewRedis()
	conn := postgres.New(c.Postgres.DataSource)
	trafficModel := model.NewSmsTrafficModel(conn, cacheRedis)
	huaweiModel := model.NewHuaweiModel(c.Huawei)
	limitRedis := c.Limiter.Redis.NewRedis()
	limiter := limit.NewPeriodLimit(c.Limiter.Expiry, c.Limiter.Quota, limitRedis, c.Limiter.KeyPrefix, limit.Align())
	return &ServiceContext{
		Config:       c,
		HuaweiModel:  huaweiModel,
		TrafficModel: trafficModel,
		Push:         c.Push,
		Limiter:      limiter,
	}
}
