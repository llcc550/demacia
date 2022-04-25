package logic

import (
	"context"
	"demacia/datacenter/position/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"demacia/datacenter/position/rpc/internal/svc"
	"demacia/datacenter/position/rpc/position"

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

	positionInfo, err := l.svcCtx.PositionModel.SelectByClassId(in.ClassId)
	if err != nil {
		return nil, status.Error(codes.NotFound, errors.PositionNotExist.Error())
	}

	return &position.PositionInfo{Id: positionInfo.Id, PositionName: positionInfo.PositionName}, nil
}
