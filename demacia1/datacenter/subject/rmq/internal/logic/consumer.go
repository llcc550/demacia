package logic

import (
	"context"
	"demacia/datacenter/class/rpc/class"
	"encoding/json"

	"demacia/common/datacenter"
	"demacia/datacenter/subject/model"
	"demacia/datacenter/subject/rmq/internal/svc"

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

func (l *Consumer) GradeConsume(_, v string) {
	var message datacenter.Message
	err := json.Unmarshal([]byte(v), &message)
	if err != nil {
		return
	}
	gradeId := message.ObjectId

	switch message.Action {
	case datacenter.Delete:
		_ = l.svcCtx.SubjectGradeModel.Deleted(&model.SubjectGrade{GradeId: gradeId})
	case datacenter.Update:
		gradeInfo, err := l.svcCtx.ClassRpc.FindGradeById(l.ctx, &class.IdReq{Id: gradeId})
		if err != nil || gradeInfo == nil {
			return
		}
		_ = l.svcCtx.SubjectGradeModel.UpdateSubjectGradeTitleByGradeId(gradeId, gradeInfo.Title)
	}
}
