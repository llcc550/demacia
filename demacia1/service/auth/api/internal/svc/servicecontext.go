package svc

import (
	"demacia/datacenter/member/rpc/memberclient"
	"demacia/datacenter/user/rpc/userclient"
	"demacia/service/auth/api/internal/config"

	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	MemberRpc memberclient.Member
	UserRpc   userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		MemberRpc: memberclient.NewMember(zrpc.MustNewClient(c.MemberRpc)),
		UserRpc:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
