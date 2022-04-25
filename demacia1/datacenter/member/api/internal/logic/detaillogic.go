package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/card/rpc/card"
	"demacia/datacenter/member/api/internal/errors"
	"demacia/datacenter/member/api/internal/svc"
	"demacia/datacenter/member/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) DetailLogic {
	return DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req types.MemberId) (resp *types.Detail, err error) {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, errlist.NoAuth
	}

	info, err := l.svcCtx.MemberModel.Details(req.MemberId, orgId)
	if err != nil {
		return nil, errors.NotExist
	}
	cards := make([]string, 0)
	teacherCardList, err := l.svcCtx.CardRpc.GetTeacherCardList(l.ctx, &card.ListReq{
		OrgId:    orgId,
		ObjectId: req.MemberId,
	})
	if err == nil && len(teacherCardList.CardNum) > 0 {
		cards = teacherCardList.CardNum
	}
	detail := types.Detail{
		MemberId: info.Id,
		UserName: info.UserName,
		TrueName: info.TrueName,
		Mobile:   info.Mobile,
		Sex:      info.Sex,
		Cards:    cards,
		Avatar:   info.Avatar,
		Face:     info.Face,
		Status:   info.Status,
	}
	return &detail, nil
}
