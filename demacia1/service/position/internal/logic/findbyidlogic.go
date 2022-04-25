package logic

import (
	"context"

	"demacia/service/position/internal/svc"
	"demacia/service/position/position"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type FindByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdLogic {
	return &FindByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByIdLogic) FindById(in *position.IdReq) (*position.PositionInfo, error) {
	// todo: add your logic here and delete this line

	return &position.PositionInfo{}, nil
}
