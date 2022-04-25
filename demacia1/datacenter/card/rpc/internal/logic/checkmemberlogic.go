package logic

import (
	"context"

	"demacia/datacenter/card/rpc/card"
	"demacia/datacenter/card/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CheckMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckMemberLogic {
	return &CheckMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckMemberLogic) CheckMember(in *card.CheckReq) (*card.BoolResp, error) {
	err := l.svcCtx.CardModel.Check(in.OrgId, in.ObjectId, roleTeacher, in.CardNum)
	return &card.BoolResp{Result: err == nil}, nil
}
