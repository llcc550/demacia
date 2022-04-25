package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/databus/rpc/databus"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"

	"demacia/datacenter/position/api/internal/svc"
	"demacia/datacenter/position/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type PositionRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPositionRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) PositionRemoveLogic {
	return PositionRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PositionRemoveLogic) PositionRemove(req types.PositionIdReq) (*types.SuccessReply, error) {

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return &types.SuccessReply{Success: false}, errlist.NoAuth
	}

	if req.Id == 0 {
		return &types.SuccessReply{Success: false}, errlist.InvalidParam
	}

	err = l.svcCtx.PositionDeviceModel.DeleteByPid(req.Id)
	if err != nil {
		return &types.SuccessReply{Success: false}, errlist.Unknown
	}

	err = l.svcCtx.PositionModel.DeleteById(oid, req.Id)
	if err != nil {
		return &types.SuccessReply{Success: false}, errlist.Unknown
	}

	threading.GoSafe(func() {
		_, _ = l.svcCtx.DataBusRpc.Delete(context.Background(), &databus.Req{Topic: datacenter.Position, ObjectId: req.Id})
	})

	return &types.SuccessReply{Success: true}, nil
}
