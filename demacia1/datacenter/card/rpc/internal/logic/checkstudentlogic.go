package logic

import (
	"context"

	"demacia/datacenter/card/rpc/card"
	"demacia/datacenter/card/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CheckStudentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckStudentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckStudentLogic {
	return &CheckStudentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckStudentLogic) CheckStudent(in *card.CheckReq) (*card.BoolResp, error) {
	err := l.svcCtx.CardModel.Check(in.OrgId, in.ObjectId, roleStudent, in.CardNum)
	return &card.BoolResp{Result: err == nil}, nil

}
