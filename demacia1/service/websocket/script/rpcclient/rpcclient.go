package main

import (
	"context"
	"flag"
	"time"

	"demacia/service/websocket/rpc/websocketclient"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

func main() {
	flag.Parse()

	client, err := zrpc.NewClient(zrpc.NewDirectClientConf([]string{"127.0.0.1:2002"}, "app", "token"))
	if err != nil {
		logx.Error(err)
	}

	cli := websocketclient.NewWebsocket(client)
	res, err := cli.Push(context.Background(), &websocketclient.Request{
		Key:  "2b1b9290-21bd-49fb-8563-068d2448aa4e",
		Code: 0,
		Msg:  "test",
	})
	if err != nil {
		logx.Error(err)
	} else {
		logx.Info(res)
	}
	time.Sleep(time.Millisecond * 100)
}
