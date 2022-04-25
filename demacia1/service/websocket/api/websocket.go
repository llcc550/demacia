package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"demacia/service/websocket/api/internal/config"
	"demacia/service/websocket/api/internal/handler"
	"demacia/service/websocket/api/internal/svc"
	"demacia/service/websocket/utils"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
)

var configFile = flag.String("f", "etc/websocket-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	redis := c.CacheRedis.NewRedis()
	serveNodeList, err := redis.Lrange(utils.WebsocketServerNodeList, 0, -1)
	if err != nil {
		serveNodeList = []string{}
	}
	serverNodesL := make([]utils.NodeInfo, 0, len(serveNodeList))
	for _, serverNode := range serveNodeList {
		var serverNodeInfo utils.NodeInfo
		_ = json.Unmarshal([]byte(serverNode), &serverNodeInfo)
		serverNodesL = append(serverNodesL, serverNodeInfo)
	}
	ctx := svc.NewServiceContext(c, serverNodesL)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
