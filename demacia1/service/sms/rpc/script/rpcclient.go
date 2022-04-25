package main

import (
	"context"
	"demacia/service/sms/model"
	"demacia/service/sms/rpc/internal/config"
	"demacia/service/sms/rpc/internal/logic"
	"demacia/service/sms/rpc/sms"
	"encoding/json"
	"flag"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
	"time"
)

var configFile = flag.String("f", "../etc/sms.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	client, err := zrpc.NewClient(zrpc.NewDirectClientConf([]string{c.ListenOn}, "app", "token"))
	if err != nil {
		logx.Error(err)
	}
	conn := client.Conn()
	cli := sms.NewSmsHandlerClient(conn)

	var mobiles []*model.ContentMobile
	mobiles = append(mobiles, &model.ContentMobile{
		TemplateId: "captcha",
		Mobile:     []string{"17551836394"},
		Params:     []string{"666"},
	})

	var req = logic.MultiBatchSendSmsRequest{
		TemplateId: "captcha",
		MultiMt:    mobiles,
	}
	bytes, err := json.Marshal(req)
	fmt.Println(string(bytes))
	response, err := cli.Push(context.Background(), &sms.RmqData{
		Data: string(bytes),
	})
	fmt.Println(err)
	logx.Info(response)
	time.Sleep(time.Millisecond * 100)
}
