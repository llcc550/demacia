package logic

import (
	"context"
	"encoding/json"

	"demacia/common/datacenter"
	"demacia/datacenter/class/rmq/internal/svc"

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

func (l *Consumer) OrganizationConsume(_, v string) {
	var message datacenter.Message
	err := json.Unmarshal([]byte(v), &message)
	if err != nil {
		return
	}
	switch message.Action {
	case datacenter.Delete:
		_ = l.svcCtx.ClassModel.DeleteByOrgId(message.ObjectId)
	}
	var data datacenter.OrganizationData
	_ = json.Unmarshal([]byte(message.Data), &data)
}
