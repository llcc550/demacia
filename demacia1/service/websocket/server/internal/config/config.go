package config

import (
	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
)

type (
	Config struct {
		ListenOn string
		Addr     string
		Tube     string
		kq.KqConf
		CacheRedis redis.RedisConf
	}
)
