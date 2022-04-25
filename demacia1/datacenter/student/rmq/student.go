package main

import (
	"context"
	"flag"
	"fmt"

	"demacia/common/datacenter"
	"demacia/datacenter/student/rmq/internal/config"
	"demacia/datacenter/student/rmq/internal/logic"
	"demacia/datacenter/student/rmq/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/student.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	l := logic.NewConsumerLogic(context.Background(), ctx)
	fmt.Printf("Starting student rmq consume\n")
	kqConf := c.KqConf
	kqConf.Topic = datacenter.Class
	q := kq.MustNewQueue(kqConf, kq.WithHandle(func(k, v string) error {
		l.ClassConsume(k, v)
		return nil
	}))
	defer q.Stop()
	fmt.Printf("Starting listening %s \n", datacenter.Class)
	q.Start()
}
