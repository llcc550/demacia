package main

import (
	"context"
	"flag"
	"time"

	"demacia/datacenter/class/rpc/classclient"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

func main() {
	flag.Parse()

	client, err := zrpc.NewClient(zrpc.NewDirectClientConf([]string{"127.0.0.1:32004"}, "app", "token"))
	if err != nil {
		logx.Error(err)
	}

	cli := classclient.NewClass(client)
	res, err := cli.ChangeStudentNum(context.Background(), &classclient.StudentNumReq{
		ClassId:    1,
		StudentNum: 100,
	})
	if err != nil {
		logx.Error(err.Error())
	} else {
		logx.Info(res)
	}
	time.Sleep(time.Second)
}
