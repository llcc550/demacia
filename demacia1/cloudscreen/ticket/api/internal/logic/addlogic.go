package logic

import (
	"context"

	"demacia/cloudscreen/ticket/api/internal/svc"
	"demacia/cloudscreen/ticket/api/internal/types"
	"demacia/cloudscreen/ticket/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddLogic {
	return AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req types.AddReq) error {
	data := make(model.Tickets, 0, len(req.OrgId))
	for _, orgId := range req.OrgId {
		data = append(data, &model.Ticket{
			Id:         0,
			UserId:     req.UserId,
			OrgId:      orgId,
			TicketDate: req.TicketDate,
		})
	}
	_ = l.svcCtx.TicketModel.BatchInsert(data)
	return nil
}
