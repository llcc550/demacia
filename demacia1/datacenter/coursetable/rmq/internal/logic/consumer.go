package logic

import (
	"context"
	"demacia/common/datacenter"
	"demacia/datacenter/coursetable/rmq/internal/svc"
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
	case datacenter.Position:
		l.positionConsume(&message)
	case datacenter.Organization:
		l.organizationConsume(&message)
	case datacenter.Member:
		l.memberConsume(&message)
	}
}

func (l *Consumer) positionConsume(message *datacenter.Message) {
	switch message.Action {
	case datacenter.Delete:
		_ = l.svcCtx.CourseTableModel.UnbindPosition(message.ObjectId)
	}
}

func (l *Consumer) organizationConsume(message *datacenter.Message) {
	switch message.Action {
	case datacenter.Delete:
		_ = l.svcCtx.CourseTableModel.DeleteByOrgId(message.ObjectId)
		_ = l.svcCtx.CourseTableDeployModel.DeleteByOrgId(message.ObjectId)
	}
}

func (l *Consumer) memberConsume(message *datacenter.Message) {
	switch message.Action {
	case datacenter.Delete:

	}
}
