package logic

import (
	"context"
	"demacia/service/sms/rpc/internal/svc"
	"demacia/service/sms/rpc/sms"
	"encoding/json"

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

func (l *PushLogic) Push(in *sms.RmqData) (*sms.Null, error) {

	if push, ok := l.svcCtx.PushMap[l.svcCtx.Config.Topic]; ok {
		s, _ := json.Marshal(&in)
		err := push.Push(string(s))
		return &sms.Null{}, err
	}

	return &sms.Null{}, nil
}
