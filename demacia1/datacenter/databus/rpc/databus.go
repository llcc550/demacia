package main

import (
	"flag"
	"fmt"

	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/databus/rpc/internal/config"
	"demacia/datacenter/databus/rpc/internal/server"
	"demacia/datacenter/databus/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/service"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/databus.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewDatabusServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		databus.RegisterDatabusServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
