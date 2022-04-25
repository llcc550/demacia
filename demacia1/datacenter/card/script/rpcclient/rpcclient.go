package main

import (
	"context"
	"flag"
	"time"

	"demacia/datacenter/card/rpc/cardclient"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

func main() {
	flag.Parse()

	client, err := zrpc.NewClient(zrpc.NewDirectClientConf([]string{"127.0.0.1:2100"}, "app", "token"))
	if err != nil {
		logx.Error(err)
	}

	cli := cardclient.NewCard(client)
	res, err := cli.GetTeacherCardList(context.Background(), &cardclient.ListReq{
		OrgId:    2,
		ObjectId: 3,
	})
	if err != nil {
		logx.Error(err)
	} else {
		logx.Info(res)
	}
	time.Sleep(time.Second)
}
