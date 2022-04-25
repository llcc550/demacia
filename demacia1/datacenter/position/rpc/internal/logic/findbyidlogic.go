package logic

import (
	"context"
	"demacia/common/errlist"
	"demacia/datacenter/position/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"demacia/datacenter/position/rpc/internal/svc"
	"demacia/datacenter/position/rpc/position"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type FindByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByPositionIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdLogic {
	return &FindByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByIdLogic) FindById(in *position.PositionIdReq) (*position.PositionInfo, error) {
	if in.PositionId == 0 {
		return nil, status.Error(codes.InvalidArgument, errlist.InvalidParam.Error())
	}

	positionInfo, err := l.svcCtx.PositionModel.SelectById(in.PositionId)
	if err != nil {
		return nil, status.Error(codes.NotFound, errors.PositionNotExist.Error())
	}

	return &position.PositionInfo{Id: positionInfo.Id, PositionName: positionInfo.PositionName}, nil
}
