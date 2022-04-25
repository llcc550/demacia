package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"demacia/service/websocket/rpc/internal/config"
	"demacia/service/websocket/rpc/internal/server"
	"demacia/service/websocket/rpc/internal/svc"
	"demacia/service/websocket/rpc/websocket"
	"demacia/service/websocket/utils"

	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/service"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/websocket.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	redis := c.CacheRedis.NewRedis()
	serveNodeList, err := redis.Lrange(utils.WebsocketServerNodeList, 0, -1)
	if err != nil || len(serveNodeList) == 0 {
		log.Fatal(err)
		return
	}
	pushMap := map[string]*kq.Pusher{}

	for _, serverNode := range serveNodeList {
		var serverNodeInfo utils.NodeInfo
		_ = json.Unmarshal([]byte(serverNode), &serverNodeInfo)
		if serverNodeInfo.Tube == "" {
			continue
		}
		pushMap[serverNodeInfo.Tube] = kq.NewPusher(c.Brokers, serverNodeInfo.Tube)
	}
	ctx := svc.NewServiceContext(c, redis, pushMap)
	srv := server.NewWebsocketServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		websocket.RegisterWebsocketServer(grpcServer, srv)
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
