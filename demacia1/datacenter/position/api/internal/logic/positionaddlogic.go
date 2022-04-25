package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/device/rpc/device"
	"demacia/datacenter/position/api/internal/svc"
	"demacia/datacenter/position/api/internal/types"
	"demacia/datacenter/position/errors"
	"demacia/datacenter/position/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type PositionAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPositionAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) PositionAddLogic {
	return PositionAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PositionAddLogic) PositionAdd(req types.PositionAddReq) (*types.SuccessReply, error) {

	if req.PositionName == "" {
		return &types.SuccessReply{Success: false}, errlist.InvalidParam
	}

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return &types.SuccessReply{Success: false}, errlist.NoAuth
	}

	positionVal, err := l.svcCtx.PositionModel.SelectByClassFullName(req.PositionName, oid)
	if positionVal.Id != 0 {
		return &types.SuccessReply{Success: false}, errors.PositionExistErr
	}

	pid, err := l.svcCtx.PositionModel.InsertPosition(&model.Position{
		Oid:          oid,
		PositionName: req.PositionName,
	})
	if err != nil {
		return &types.SuccessReply{Success: false}, errlist.Unknown
	}

	if len(req.DeviceIds) != 0 {
		psVal, err := l.svcCtx.PositionDeviceModel.SelectByDeviceIds(req.DeviceIds)
		if err == sql.ErrNoRows {
		} else {
			if len(psVal) > 0 {
				return &types.SuccessReply{Success: true}, errors.PositionBindDeviceErr
			}
		}
		devices := make([]*model.PositionDevice, 0, len(req.DeviceIds))
		for _, sid := range req.DeviceIds {
			deviceInfo, err := l.svcCtx.DeviceRpc.GetDeviceInfoById(l.ctx, &device.IdReq{Id: sid})
			if err != nil {
				l.Logger.Errorf("call deviceRpc error:%s", err.Error())
				continue
			}
			devices = append(devices, &model.PositionDevice{
				PositionId: pid,
				DeviceId:   sid,
				DeviceName: deviceInfo.Title,
			})
		}
		if err := l.svcCtx.PositionDeviceModel.InsertPositionDevice(devices); err != nil {
			l.Logger.Errorf("insert position_device error:%s", err.Error())
		}
	}

	return &types.SuccessReply{Success: true}, nil
}
