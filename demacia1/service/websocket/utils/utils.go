package utils

type (
	NodeInfo struct {
		Tube string `json:"tube"`
		Addr string `json:"addr"`
	}
	WebSocketPush struct {
		Key  string `json:"key"`
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
	}
)

const (
	WebsocketConnToServerNodePrefix = "cache:websocket:conn-to-server-node:"
	WebsocketServerNodeList         = "cache:websocket:server-list"
	Ttl                             = 3600
)
