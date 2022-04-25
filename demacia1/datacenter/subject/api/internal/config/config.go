package config

import (
	"demacia/common/baseauth"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Postgres struct {
		DataSource string
	}
	CacheRedis redis.RedisConf
	Auth       baseauth.AuthConfig
	ClassRpc   zrpc.RpcClientConf
	MemberRpc  zrpc.RpcClientConf
	DataBusRpc zrpc.RpcClientConf
}
