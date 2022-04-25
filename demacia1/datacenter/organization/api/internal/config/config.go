package config

import (
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
	}
	Postgres struct {
		DataSource string
	}
	CacheRedis redis.RedisConf
	CommonRpc  zrpc.RpcClientConf
	MemberRpc  zrpc.RpcClientConf
	DataBusRpc zrpc.RpcClientConf
}
