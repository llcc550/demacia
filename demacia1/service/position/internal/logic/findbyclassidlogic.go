package logic

import (
	"context"

	"demacia/service/position/internal/svc"
	"demacia/service/position/position"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type FindByClassIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByClassIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByClassIdLogic {
	return &FindByClassIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByClassIdLogic) FindByClassId(in *position.ClassIdReq) (*position.PositionInfo, error) {
	// todo: add your logic here and delete this line

	return &position.PositionInfo{}, nil
}
