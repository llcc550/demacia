package config

import (
	"demacia/common/baseauth"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Postgres struct {
		DataSource string
	}
	CacheRedis   redis.RedisConf
	Auth         baseauth.AuthConfig
	ClassRpc     zrpc.RpcClientConf
	StudentRpc   zrpc.RpcClientConf
	WebsocketRpc zrpc.RpcClientConf
	Brokers      []string
}
