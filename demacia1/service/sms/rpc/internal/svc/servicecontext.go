package svc

import (
	"demacia/service/sms/model"
	"demacia/service/sms/rpc/internal/config"
	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/limit"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"regexp"
)

type ServiceContext struct {
	Config        config.Config
	contentRegexp *regexp.Regexp
	TrafficModel  *model.SmsTrafficModel
	HuaweiModel   *model.HuaweiModel
	Limiter       *limit.PeriodLimit
	cacheRedis    *redis.Redis
	templates     map[string][]string
	Push          bool
	PushMap       map[string]*kq.Pusher
}

func NewServiceContext(c config.Config, pushMap map[string]*kq.Pusher) *ServiceContext {
	cacheRedis := c.CacheRedis.NewRedis()
	conn := postgres.New(c.Postgres.DataSource)
	trafficModel := model.NewSmsTrafficModel(conn, cacheRedis)
	huaweiModel := model.NewHuaweiModel(c.Huawei)
	limitRedis := c.Limiter.Redis.NewRedis()
	limiter := limit.NewPeriodLimit(c.Limiter.Expiry, c.Limiter.Quota, limitRedis, c.Limiter.KeyPrefix, limit.Align())
	return &ServiceContext{
		Config:       c,
		TrafficModel: trafficModel,
		HuaweiModel:  huaweiModel,
		Limiter:      limiter,
		Push:         c.Push,
		cacheRedis:   cacheRedis,
		PushMap:      pushMap,
	}
}
