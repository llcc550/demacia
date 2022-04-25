package logic

import (
	"context"
	"encoding/json"

	"demacia/common/datacenter"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/databus/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *databus.Req) (*databus.Res, error) {
	s, _ := json.Marshal(datacenter.Message{
		Topic:    in.Topic,
		ObjectId: in.ObjectId,
		Action:   datacenter.Update,
	})
	_ = l.svcCtx.KqPusher.Push(string(s))
	return &databus.Res{Result: true}, nil
}
