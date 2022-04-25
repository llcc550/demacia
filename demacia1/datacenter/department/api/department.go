package main

import (
	"flag"
	"fmt"

	"demacia/datacenter/department/api/internal/config"
	"demacia/datacenter/department/api/internal/handler"
	"demacia/datacenter/department/api/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
)

var configFile = flag.String("f", "etc/department-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
