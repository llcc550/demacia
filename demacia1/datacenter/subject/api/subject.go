package main

import (
	"demacia/datacenter/subject/api/internal/config"
	"demacia/datacenter/subject/api/internal/handler"
	"demacia/datacenter/subject/api/internal/svc"
	"flag"
	"fmt"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
)

var configFile = flag.String("f", "etc/subject-api.yaml", "the config file")

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
