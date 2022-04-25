package main

import (
	"context"
	"demacia/common/datacenter"
	"demacia/datacenter/coursetable/rmq/internal/config"
	"demacia/datacenter/coursetable/rmq/internal/logic"
	"demacia/datacenter/coursetable/rmq/internal/svc"
	"flag"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/coursetable.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	l := logic.NewConsumerLogic(context.Background(), ctx)

	kqConf := c.KqConf
	kqConf.Group = datacenter.CourseTable
	kqConf.Topic = datacenter.Kafka
	q := kq.MustNewQueue(kqConf, kq.WithHandle(func(k, v string) error {
		l.Consume(k, v)
		return nil
	}))
	defer q.Stop()
	fmt.Printf("Starting coursetable rmq consume\n")
	q.Start()
}
