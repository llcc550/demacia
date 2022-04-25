package logic

import (
	"context"

	"demacia/common/errlist"
	"demacia/datacenter/member/api/internal/svc"
	"demacia/datacenter/member/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type StatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) StatusLogic {
	return StatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatusLogic) Status(req types.StatusReq) error {
	if req.Status != 1 && req.Status != -1 {
		return errlist.InvalidParam
	}
	return l.svcCtx.MemberModel.SetMemberStatus(req.Status, req.MemberId)
}
