package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"

	"demacia/service/event/api/internal/svc"
	"demacia/service/event/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type EventDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEventDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) EventDeleteLogic {
	return EventDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EventDeleteLogic) EventDelete(req types.EventIdReq) error {
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
	err = l.svcCtx.EventModel.Delete(req.EventId)
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		_ = l.svcCtx.EventPositionModel.Delete(req.EventId)
	})
	return nil
}
