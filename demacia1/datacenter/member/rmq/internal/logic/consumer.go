package logic

import (
	"context"
	"encoding/json"

	"demacia/common/datacenter"
	"demacia/datacenter/member/rmq/internal/svc"

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
	case datacenter.Organization:
		l.organizationConsume(&message)
	}
}

func (l *Consumer) organizationConsume(message *datacenter.Message) {
	switch message.Action {
	case datacenter.Delete:
		_ = l.svcCtx.MemberModel.DeleteByOrgId(message.ObjectId)
	}
}
