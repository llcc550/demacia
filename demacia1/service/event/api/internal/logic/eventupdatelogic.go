package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/service/event/api/internal/svc"
	"demacia/service/event/api/internal/types"
	"demacia/service/event/model"
	"demacia/service/position/position"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
	"time"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type EventUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEventUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) EventUpdateLogic {
	return EventUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EventUpdateLogic) EventUpdate(req types.EventInsertReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	eventInfo, err := l.svcCtx.EventModel.FindOne(req.EventId)
	if err != nil {
		return errlist.EventNotExist
	}
	if eventInfo.OrgId != orgId {
		return errlist.NoAuth
	}
	if req.EndTime <= req.StartTime || req.EndTime < time.Now().Unix() {
		return errlist.EventTimeErr
	}
	if len(req.PositionIds) < 1 {
		return errlist.EventPositionErr
	}
	err = l.svcCtx.EventModel.Update(&model.Event{
		Id:         req.EventId,
		OrgId:      orgId,
		Name:       req.EventName,
		CategoryId: req.CategoryId,
		Content:    req.EventContent,
		PushType:   req.PushType,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	})
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		for _, id := range req.PositionIds {
			positionInfo, err := l.svcCtx.PositionRpc.FindById(context.Background(), &position.IdReq{Id: id})
			if err != nil {
				logx.Errorf("get positionInfo err:%s", err.Error())
				continue
			}
			err = l.svcCtx.EventPositionModel.Delete(req.EventId)
			if err != nil {
				continue
			}
			err = l.svcCtx.EventPositionModel.Insert(&model.EventPosition{
				EventId:       req.EventId,
				PositionId:    positionInfo.Id,
				PositionTitle: positionInfo.PositionName,
			})
			if err != nil {
				logx.Errorf("insert into event_range err:%s", err.Error())
				continue
			}
		}

	})
	return nil
}
