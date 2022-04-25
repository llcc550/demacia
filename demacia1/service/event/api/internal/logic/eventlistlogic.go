package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/service/event/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"

	"demacia/service/event/api/internal/svc"
	"demacia/service/event/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type EventListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEventListLogic(ctx context.Context, svcCtx *svc.ServiceContext) EventListLogic {
	return EventListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EventListLogic) EventList(req types.ListReq) (resp *types.ListResponse, err error) {
	// todo: add your logic here and delete this line
	resp = &types.ListResponse{
		List:  []*types.EventDetail{},
		Count: 0,
	}
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, err
	}
	//根据position_ids获取event_ids
	eventIds := make([]int64, 0)
	if len(req.PositionIds) != 0 {
		eventRangeList, err := l.svcCtx.EventPositionModel.FindListByPositionIds(req.PositionIds)
		if err == nil {
			for _, v := range eventRangeList {
				eventIds = append(eventIds, v.EventId)
			}
		}
	}
	eventList, count, err := l.svcCtx.EventModel.FindListByConditions(eventIds, req.EventName, req.MemberName, orgId, req.StartTime, req.EndTime, req.CategoryId, req.Page, req.Limit)
	if err != nil {
		return resp, err
	}
	//获取category Map
	categoryMap := map[int64]model.Category{}
	eventMap := map[int64][]*model.EventPosition{}

	mr.FinishVoid(func() {
		categoryList, err := l.svcCtx.CategoryModel.FindList(orgId, "")
		if err == nil {
			for _, v := range *categoryList {
				categoryMap[v.Id] = model.Category{
					Id:    v.Id,
					OrgId: v.OrgId,
					Name:  v.Name,
				}
			}
		}
	}, func() {
		//获取所有的position
		ids := make([]int64, 0)
		for _, v := range eventList {
			ids = append(ids, v.Id)
		}
		positionList, err := l.svcCtx.EventPositionModel.FindListByEventIds(ids)
		if err == nil {
			for _, v := range positionList {
				eventMap[v.EventId] = append(eventMap[v.EventId], &model.EventPosition{
					Id:            v.Id,
					EventId:       v.EventId,
					PositionId:    v.PositionId,
					PositionTitle: v.PositionTitle,
				})
			}
		}
	})
	//组数据
	for _, v := range eventList {
		categoryName := ""
		position := make([]*types.Position, 0)
		if category, ok := categoryMap[v.CategoryId]; ok {
			categoryName = category.Name
		}
		if positions, ok := eventMap[v.Id]; ok {
			for _, eventPosition := range positions {
				position = append(position, &types.Position{
					PositionId:    eventPosition.PositionId,
					PositionTitle: eventPosition.PositionTitle,
				})
			}

		}
		resp.List = append(resp.List, &types.EventDetail{
			EventId:      v.Id,
			EventName:    v.Name,
			EventContent: v.Content,
			CategoryId:   v.CategoryId,
			CategoryName: categoryName,
			PushType:     v.PushType,
			StartTime:    v.StartTime,
			EndTime:      v.EndTime,
			CreatedAt:    v.CreatedAt,
			MemberId:     v.MemberId,
			MemberName:   v.MemberName,
			Position:     position,
		})

	}
	resp.Count = count
	return resp, nil
}
