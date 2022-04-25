package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/member/api/internal/svc"
	"demacia/datacenter/member/api/internal/types"
	"demacia/datacenter/member/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListLogic {
	return ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req types.ListReq) (resp *types.ListRes, err error) {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, errlist.NoAuth
	}
	count, list, err := l.svcCtx.MemberModel.AdvancedSearchList(model.ListCondition{
		OrgId:      orgId,
		Page:       req.Page,
		Limit:      req.Limit,
		TrueName:   req.TrueName,
		FaceStatus: req.FaceStatus,
	})
	res := types.ListRes{List: []types.Member{}, Count: count}
	if err != nil || count == 0 {
		return &res, err
	}
	for _, v := range list {
		res.List = append(res.List, types.Member{
			MemberId:   v.Id,
			UserName:   v.UserName,
			TrueName:   v.TrueName,
			Mobile:     v.Mobile,
			Status:     v.Status,
			Face:       v.Face,
			FaceStatus: v.FaceStatus,
			Avatar:     v.Avatar,
			JoinTime:   v.JoinDate,
		})
	}
	return &res, nil
}
