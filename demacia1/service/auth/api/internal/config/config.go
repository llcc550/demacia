package config

import (
	"demacia/common/baseauth"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth      baseauth.AuthConfig
	MemberRpc zrpc.RpcClientConf
	UserRpc   zrpc.RpcClientConf
}
