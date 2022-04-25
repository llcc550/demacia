package main

import (
	"context"
	"flag"
	"fmt"

	"demacia/common/datacenter"
	"demacia/datacenter/class/rmq/internal/config"
	"demacia/datacenter/class/rmq/internal/logic"
	"demacia/datacenter/class/rmq/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
)

var configFile = flag.String("f", "etc/class.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	l := logic.NewConsumerLogic(context.Background(), ctx)
	fmt.Printf("Starting class rmq consume\n")
	mr.FinishVoid(func() {
		kqConf := c.KqConf
		kqConf.Topic = datacenter.Organization
		q := kq.MustNewQueue(kqConf, kq.WithHandle(func(k, v string) error {
			l.OrganizationConsume(k, v)
			return nil
		}))
		defer q.Stop()
		fmt.Printf("Starting listening %s \n", datacenter.Organization)
		q.Start()
	})
}
