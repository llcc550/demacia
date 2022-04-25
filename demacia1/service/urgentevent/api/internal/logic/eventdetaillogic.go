package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/member/rpc/member"
	"demacia/service/urgentevent/api/internal/svc"
	"demacia/service/urgentevent/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
)

type EventDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEventDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) EventDetailLogic {
	return EventDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EventDetailLogic) EventDetail(req types.EventIdReq) (resp *types.EventDetail, err error) {
	resp = &types.EventDetail{
		EventId:      0,
		EventName:    "",
		EventContent: "",
		CategoryId:   0,
		CategoryName: "",
		PushType:     0,
		StartTime:    0,
		EndTime:      0,
		CreatedAt:    0,
		MemberId:     0,
		MemberName:   "",
		Position:     []*types.Position{},
	}
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, err
	}
	urgenteventInfo, err := l.svcCtx.EventModel.FindOne(req.EventId)
	if err != nil {
		return resp, errlist.EventNotExist
	}
	if urgenteventInfo.OrgId != orgId {
		return resp, errlist.NoAuth
	}
	categoryName := ""
	memberName := ""
	position := make([]*types.Position, 0)
	_ = mr.Finish(func() error {
		categoryInfo, err := l.svcCtx.CategoryModel.FindOne(urgenteventInfo.CategoryId)
		if err == nil {
			categoryName = categoryInfo.Name
		}
		return nil
	}, func() error {
		memberInfo, err := l.svcCtx.MemberRpc.FindOneById(l.ctx, &member.IdReq{Id: urgenteventInfo.MemberId})
		if err == nil {
			memberName = memberInfo.TrueName
		}
		return nil
	}, func() error {
		urgenteventRange, err := l.svcCtx.EventPositionModel.FindListByEventIds([]int64{req.EventId})
		if err == nil {
			for _, v := range urgenteventRange {
				position = append(position, &types.Position{
					PositionId:    v.PositionId,
					PositionTitle: v.PositionTitle,
				})
			}
		}
		return nil
	})
	resp.EventId = urgenteventInfo.Id
	resp.EventName = urgenteventInfo.Name
	resp.EventContent = urgenteventInfo.Content
	resp.CategoryId = urgenteventInfo.CategoryId
	resp.CategoryName = categoryName
	resp.PushType = urgenteventInfo.PushType
	resp.StartTime = urgenteventInfo.StartTime
	resp.EndTime = urgenteventInfo.EndTime
	resp.CreatedAt = urgenteventInfo.CreatedAt
	resp.MemberId = urgenteventInfo.MemberId
	resp.MemberName = memberName
	resp.Position = position
	return resp, nil
}
