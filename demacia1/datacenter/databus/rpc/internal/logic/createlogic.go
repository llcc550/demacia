package logic

import (
	"context"
	"encoding/json"

	"demacia/common/datacenter"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/databus/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *databus.Req) (*databus.Res, error) {
	s, _ := json.Marshal(datacenter.Message{
		Topic:    in.Topic,
		ObjectId: in.ObjectId,
		Action:   datacenter.Create,
	})
	_ = l.svcCtx.KqPusher.Push(string(s))
	return &databus.Res{Result: true}, nil
}
