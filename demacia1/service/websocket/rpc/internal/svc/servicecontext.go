package svc

import (
	"demacia/service/websocket/rpc/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config  config.Config
	Redis   *redis.Redis
	PushMap map[string]*kq.Pusher
}

func NewServiceContext(c config.Config, redis *redis.Redis, pushMap map[string]*kq.Pusher) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Redis:   redis,
		PushMap: pushMap,
	}
}
