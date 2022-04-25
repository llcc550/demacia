package config

import (
	"demacia/service/sms/model"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
)

type Config struct {
	ListenOn string
	kq.KqConf
	Limiter struct {
		Expiry    int
		KeyPrefix string
		Quota     int
		Redis     redis.RedisConf
	}
	Push       bool
	CacheRedis redis.RedisConf
	Postgres   struct {
		DataSource string
	}
	Huawei *model.HuaweiConfig
}
