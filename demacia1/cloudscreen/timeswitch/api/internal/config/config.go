package config

import (
	"demacia/common/baseauth"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth     baseauth.AuthConfig
	Postgres struct {
		DataSource string
	}
	CacheRedis redis.RedisConf
	CommonRpc  zrpc.RpcClientConf
	DataBusRpc zrpc.RpcClientConf
}
