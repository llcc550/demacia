package config

import (
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Postgres struct {
		DataSource string
	}
	CacheRedis redis.RedisConf
}
