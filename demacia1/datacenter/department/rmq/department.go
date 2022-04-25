package main

import (
	"context"
	"flag"
	"fmt"

	"demacia/common/datacenter"
	"demacia/datacenter/department/rmq/internal/config"
	"demacia/datacenter/department/rmq/internal/logic"
	"demacia/datacenter/department/rmq/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
)

var configFile = flag.String("f", "etc/department.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	l := logic.NewConsumerLogic(context.Background(), ctx)
	fmt.Printf("Starting department rmq consume\n")

	mr.FinishVoid(func() {
		kqConf := c.KqConf
		kqConf.Topic = datacenter.Member
		q := kq.MustNewQueue(kqConf, kq.WithHandle(func(k, v string) error {
			l.MemberConsume(k, v)
			return nil
		}))
		defer q.Stop()
		fmt.Printf("Starting listening %s \n", datacenter.Member)
		q.Start()
	})
}
