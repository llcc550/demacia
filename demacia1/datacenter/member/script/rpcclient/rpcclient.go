package main

import (
	"context"
	"flag"
	"time"

	"demacia/datacenter/member/rpc/memberclient"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

func main() {
	flag.Parse()
	client, err := zrpc.NewClient(zrpc.NewDirectClientConf([]string{"127.0.0.1:2100"}, "app", "token"))
	if err != nil {
		logx.Error(err)
	}
	cli := memberclient.NewMember(client)

	// test FindOneByUserName
	res1, err1 := cli.FindOneByUserName(context.Background(), &memberclient.UserNameReq{
		UserName: "little bitch",
	})
	if err1 != nil {
		logx.Error(err1)
	} else {
		logx.Info(res1)
	}
	time.Sleep(time.Millisecond * 100)

	// test Insert
	res2, err2 := cli.Insert(context.Background(), &memberclient.InsertReq{
		UserName: "little bitch",
		Mobile:   "13812670001",
		TrueName: "little bitch ",
		Password: "little bitch's password",
		OrgId:    2,
	})
	if err2 != nil {
		logx.Error(err2)
	} else {
		logx.Info(res2)
	}
	time.Sleep(time.Millisecond * 100)
}
