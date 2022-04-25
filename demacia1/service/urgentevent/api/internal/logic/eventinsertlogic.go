package logic

import (
	"context"
	"time"

	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/member/rpc/memberclient"
	"demacia/service/position/position"
	"demacia/service/urgentevent/api/internal/svc"
	"demacia/service/urgentevent/api/internal/types"
	"demacia/service/urgentevent/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type EventInsertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEventInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) EventInsertLogic {
	return EventInsertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EventInsertLogic) EventInsert(req types.EventInsertReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	memberId, err := baseauth.GetMemberId(l.ctx)
	if err != nil {
		return err
	}
	memberInfo, err := l.svcCtx.MemberRpc.FindOneById(l.ctx, &memberclient.IdReq{Id: memberId})
	if err != nil {
		logx.Errorf("get memberInfo err:%s", err.Error())
		return errlist.MemberNotExist
	}
	if len(req.PositionIds) < 1 {
		return errlist.EventPositionErr
	}
	if req.EndTime <= req.StartTime || req.EndTime < time.Now().Unix() {
		return errlist.EventTimeErr
	}
	urgenteventId, err := l.svcCtx.EventModel.Insert(&model.Event{
		Id:         0,
		OrgId:      orgId,
		Name:       req.EventName,
		CategoryId: req.CategoryId,
		Content:    req.EventContent,
		PushType:   req.PushType,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		MemberId:   memberId,
		MemberName: memberInfo.TrueName,
		CreatedAt:  time.Now().Unix(),
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
			err = l.svcCtx.EventPositionModel.Insert(&model.EventPosition{
				EventId:       urgenteventId,
				PositionId:    positionInfo.Id,
				PositionTitle: positionInfo.PositionName,
			})
			if err != nil {
				logx.Errorf("insert into urgentevent_range err:%s", err.Error())
				continue
			}
		}

	})
	return nil
}
