package config

import (
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type Config struct {
	Postgres struct {
		DataSource string
	}
	CacheRedis     redis.RedisConf
	CourseTableRpc zrpc.RpcClientConf
	ClassRpc       zrpc.RpcClientConf
	MemberRpc      zrpc.RpcClientConf
}
