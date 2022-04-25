package logic

import (
	"context"
	"demacia/cloudscreen/timeswitch/rmq/internal/svc"
	"demacia/common/datacenter"
	"encoding/json"

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
	case datacenter.Device:
		l.deviceConsume(&message)
	}
}

func (l *Consumer) deviceConsume(message *datacenter.Message) {
	switch message.Action {
	case datacenter.Delete:
		_ = l.svcCtx.TimeSwitchModel.DeleteByDeviceId(message.ObjectId)
		_ = l.svcCtx.TimeSwitchDateModel.DeleteByDeviceId(message.ObjectId)
		_ = l.svcCtx.TimeSwitchConfigModel.DeleteByDeviceId(message.ObjectId)
	}
}
