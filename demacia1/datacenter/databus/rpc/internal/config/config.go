package config

import "gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Brokers []string
	Mongo   struct {
		Url        string
		Collection struct {
			Log string
		}
	}
}
