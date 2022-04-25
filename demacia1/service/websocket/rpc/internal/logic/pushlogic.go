package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"demacia/service/websocket/rpc/internal/svc"
	"demacia/service/websocket/rpc/websocket"
	"demacia/service/websocket/utils"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type PushLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushLogic {
	return &PushLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushLogic) Push(in *websocket.Request) (*websocket.Null, error) {
	serverConn := fmt.Sprintf("%s%s", utils.WebsocketConnToServerNodePrefix, in.Key)
	serverNode, err := l.svcCtx.Redis.Get(serverConn)
	if err != nil || serverNode == "" {
		return nil, err
	}
	if push, ok := l.svcCtx.PushMap[serverNode]; ok {
		data := utils.WebSocketPush{
			Key:  in.Key,
			Code: in.Code,
			Msg:  in.Msg,
		}
		s, _ := json.Marshal(&data)
		_ = push.Push(string(s))
	}
	return &websocket.Null{}, nil
}
