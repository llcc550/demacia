package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/member/rpc/member"
	"demacia/service/event/api/internal/svc"
	"demacia/service/event/api/internal/types"

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
	eventInfo, err := l.svcCtx.EventModel.FindOne(req.EventId)
	if err != nil {
		return resp, errlist.EventNotExist
	}
	if eventInfo.OrgId != orgId {
		return resp, errlist.NoAuth
	}
	categoryName := ""
	memberName := ""
	position := make([]*types.Position, 0)
	_ = mr.Finish(func() error {
		categoryInfo, err := l.svcCtx.CategoryModel.FindOne(eventInfo.CategoryId)
		if err == nil {
			categoryName = categoryInfo.Name
		}
		return nil
	}, func() error {
		memberInfo, err := l.svcCtx.MemberRpc.FindOneById(l.ctx, &member.IdReq{Id: eventInfo.MemberId})
		if err == nil {
			memberName = memberInfo.TrueName
		}
		return nil
	}, func() error {
		eventRange, err := l.svcCtx.EventPositionModel.FindListByEventIds([]int64{req.EventId})
		if err == nil {
			for _, v := range eventRange {
				position = append(position, &types.Position{
					PositionId:    v.PositionId,
					PositionTitle: v.PositionTitle,
				})
			}
		}
		return nil
	})
	resp.EventId = eventInfo.Id
	resp.EventName = eventInfo.Name
	resp.EventContent = eventInfo.Content
	resp.CategoryId = eventInfo.CategoryId
	resp.CategoryName = categoryName
	resp.PushType = eventInfo.PushType
	resp.StartTime = eventInfo.StartTime
	resp.EndTime = eventInfo.EndTime
	resp.CreatedAt = eventInfo.CreatedAt
	resp.MemberId = eventInfo.MemberId
	resp.MemberName = memberName
	resp.Position = position
	return resp, nil
}
