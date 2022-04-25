package main

import (
	"context"
	"flag"
	"time"

	"demacia/datacenter/organization/rpc/organizationclient"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	flag.Parse()
	client, err := zrpc.NewClient(zrpc.NewDirectClientConf([]string{"127.0.0.1:2100"}, "app", "token"))
	if err != nil {
		logx.Error(err)
		return
	}
	id := int64(2)
	cli := organizationclient.NewOrganization(client)
	res, err := cli.FindOne(context.Background(), &organizationclient.IdReply{
		Id: id,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			logx.Infof("organization not find. id is %d", id)
		} else {
			logx.Errorf("find one err. id is %d,err is %s", id, err.Error())
		}
	} else {
		logx.Info(res)
	}
	time.Sleep(time.Millisecond * 100)
}
