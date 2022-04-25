package config

import (
	"demacia/service/sms/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	CacheRedis redis.RedisConf
	Brokers    []string
	Topic      string
	Push       bool                `json:"Push,optional"`
	Huawei     *model.HuaweiConfig `json:"Huawei,optional"`
	Postgres   struct {
		DataSource string
	}
	Limiter struct {
		Expiry    int
		KeyPrefix string
		Quota     int
		Redis     redis.RedisConf
	}
}
