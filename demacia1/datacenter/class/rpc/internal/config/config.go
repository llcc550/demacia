package config

import (
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Postgres struct {
		DataSource string
	}
	CacheRedis redis.RedisConf
}
