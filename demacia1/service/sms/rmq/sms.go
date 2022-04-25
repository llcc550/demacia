package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"demacia/service/sms/rmq/internal/config"
	"demacia/service/sms/rmq/internal/logic"
	"demacia/service/sms/rmq/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

var configFile = flag.String("f", "etc/sms.yaml", "the config file for sms rmq service")

func main() {

	flag.Parse()

	var c config.Config

	conf.MustLoad(*configFile, &c)

	consumerLogic := logic.NewConsumerLogic(context.Background(), svc.NewServiceContext(c))

	threading.GoSafe(func() {
		q := kq.MustNewQueue(c.KqConf, kq.WithHandle(func(k, v string) error {
			err := consumerLogic.Consumer(v)
			return err
		}))
		defer q.Stop()
		q.Start()
	})

	fmt.Printf("starting sms rmq server...\n")
	_ = http.ListenAndServe(c.ListenOn, nil)
}
