package main

import (
	"flag"
	"fmt"

	"demacia/datacenter/student/api/internal/config"
	"demacia/datacenter/student/api/internal/handler"
	"demacia/datacenter/student/api/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
)

var configFile = flag.String("f", "etc/student-api.yaml", "the config file")

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

// todo: 1、中间件；
// todo: 2、导出启用websocket
// todo: 3、错误列表、readme.md、sql
