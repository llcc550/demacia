package main

import (
	"context"
	"flag"
	"fmt"

	"demacia/common/datacenter"
	"demacia/datacenter/member/rmq/internal/config"
	"demacia/datacenter/member/rmq/internal/logic"
	"demacia/datacenter/member/rmq/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/member.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	l := logic.NewConsumerLogic(context.Background(), ctx)

	kqConf := c.KqConf
	kqConf.Group = datacenter.Member
	kqConf.Topic = datacenter.Kafka
	q := kq.MustNewQueue(kqConf, kq.WithHandle(func(k, v string) error {
		l.Consume(k, v)
		return nil
	}))
	defer q.Stop()
	fmt.Printf("Starting member rmq consume\n")
	q.Start()
}
