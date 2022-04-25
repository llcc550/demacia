package main

import (
	"demacia/service/sms/rpc/internal/config"
	"demacia/service/sms/rpc/internal/server"
	"demacia/service/sms/rpc/internal/svc"
	"demacia/service/sms/rpc/sms"
	"flag"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/service"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/sms.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	pushMap := map[string]*kq.Pusher{}

	pushMap[c.Topic] = kq.NewPusher(c.Brokers, c.Topic)

	ctx := svc.NewServiceContext(c, pushMap)
	srv := server.NewSmsHandlerServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		sms.RegisterSmsHandlerServer(grpcServer, srv)
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
