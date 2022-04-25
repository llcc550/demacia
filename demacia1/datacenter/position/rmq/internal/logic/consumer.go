package logic

import (
	"context"
	"encoding/json"

	"demacia/common/datacenter"
	"demacia/datacenter/position/rmq/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type Consumer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConsumerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Consumer {
	return &Consumer{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Consumer) Consume(_, v string) {
	var message datacenter.Message
	err := json.Unmarshal([]byte(v), &message)
	if err != nil {
		return
	}
	switch message.Topic {
	case datacenter.Class:
		l.classConsume(&message)
	}
}

func (l *Consumer) classConsume(message *datacenter.Message) {
	switch message.Action {
	case datacenter.Delete:
		_ = l.svcCtx.PositionModel.UnbindClass(message.ObjectId)
	}
}
