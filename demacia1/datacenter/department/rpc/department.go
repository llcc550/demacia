package main

import (
	"flag"
	"fmt"

	"demacia/datacenter/department/rpc/department"
	"demacia/datacenter/department/rpc/internal/config"
	"demacia/datacenter/department/rpc/internal/server"
	"demacia/datacenter/department/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/service"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/department.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewDepartmentServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		department.RegisterDepartmentServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
