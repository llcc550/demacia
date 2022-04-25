package main

import (
	"context"
	"flag"
	"fmt"

	"demacia/common/datacenter"
	"demacia/datacenter/position/rmq/internal/config"
	"demacia/datacenter/position/rmq/internal/logic"
	"demacia/datacenter/position/rmq/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
)

var configFile = flag.String("f", "etc/position.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	l := logic.NewConsumerLogic(context.Background(), ctx)
	fmt.Printf("Starting position rmq consume\n")
	mr.FinishVoid(func() {
		kqConf := c.KqConf
		kqConf.Group = datacenter.Position
		kqConf.Topic = datacenter.Kafka
		q := kq.MustNewQueue(kqConf, kq.WithHandle(func(k, v string) error {
			l.Consume(k, v)
			return nil
		}))
		defer q.Stop()
		fmt.Printf("Starting coursetable rmq consume\n")
		q.Start()
	})
}
