package config

import (
	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/service"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type (
	Config struct {
		service.ServiceConf
		kq.KqConf
		Postgres struct {
			DataSource string
		}
		CacheRedis redis.RedisConf
		MemberRpc  zrpc.RpcClientConf
	}
)
