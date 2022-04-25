package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"sync"

	"demacia/service/websocket/server/internal/config"
	"demacia/service/websocket/utils"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gitlab.u-jy.cn/xiaoyang/go-queue/kq"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/conf"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type (
	wsConn struct {
		*websocket.Conn
		Mux sync.RWMutex
	}
)

var configFile = flag.String("f", "etc/websocket.yaml", "the config file")

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		fmt.Println(r.Header)
		return true
	},
}
var connMap sync.Map

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	threading.GoSafe(func() {
		q := kq.MustNewQueue(c.KqConf, kq.WithHandle(func(k, v string) error {
			return Consume(v)
		}))
		defer q.Stop()
		q.Start()
	})
	redis := c.CacheRedis.NewRedis()
	nodeInfo := utils.NodeInfo{
		Tube: c.Tube,
		Addr: c.Addr,
	}
	s, _ := json.Marshal(nodeInfo)
	_, _ = redis.Lrem(utils.WebsocketServerNodeList, 0, string(s))
	_, _ = redis.Lpush(utils.WebsocketServerNodeList, string(s))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrade.Upgrade(w, r, nil)
		if err != nil {
			logx.Error(err)
			return
		}
		defer conn.Close()
		connKey := uuid.New().String()
		serverConnInServer := fmt.Sprintf("%s%s", utils.WebsocketConnToServerNodePrefix, connKey)
		_ = redis.Set(serverConnInServer, c.Tube)
		_ = redis.Expire(serverConnInServer, utils.Ttl)
		connWithLock := wsConn{
			Conn: conn,
			Mux:  sync.RWMutex{},
		}
		connMap.Store(connKey, &connWithLock)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(connKey))
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				_, _ = redis.Del(serverConnInServer)
				break
			} else {
				_ = redis.Expire(serverConnInServer, utils.Ttl)
			}
		}
	})
	fmt.Printf("Starting websocket server at %s...\n", c.Addr)
	_ = http.ListenAndServe(c.ListenOn, nil)
}

func Consume(pushData string) error {
	logx.Infof("开始处理推送: %s", pushData)
	var data utils.WebSocketPush
	_ = json.Unmarshal([]byte(pushData), &data)
	if data.Key == "" || data.Msg == "" {
		return nil
	}
	if connKey, ok := connMap.Load(data.Key); ok {
		if c, ok := connKey.(*wsConn); ok {
			c.Mux.Lock()
			defer c.Mux.Unlock()
			s := fmt.Sprintf(`{"code":%d,"msg":"%s"}`, data.Code, data.Msg)
			_ = c.Conn.WriteMessage(websocket.TextMessage, []byte(s))
		} else {
			connMap.Delete(data.Key)
		}
	}
	return nil
}
