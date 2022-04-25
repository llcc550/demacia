package logic

import (
	"context"
	"time"

	"demacia/common/errlist"
	"demacia/service/websocket/api/internal/svc"
	"demacia/service/websocket/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type NodesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) NodesLogic {
	return NodesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NodesLogic) Nodes() (*types.Response, error) {
	if len(l.svcCtx.Nodes) == 0 {
		return nil, errlist.WebsocketNodeOffline
	}
	res := types.Response{
		Recommend: "",
		List:      []string{},
	}
	count := int64(len(l.svcCtx.Nodes))
	now := time.Now().UnixMicro()
	for k, v := range l.svcCtx.Nodes {
		if now%count == int64(k) {
			res.Recommend = v.Addr
			continue
		}
		res.List = append(res.List, v.Addr)
	}
	return &res, nil
}
