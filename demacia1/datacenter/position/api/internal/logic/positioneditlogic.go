package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/device/rpc/device"
	"demacia/datacenter/position/api/internal/svc"
	"demacia/datacenter/position/api/internal/types"
	"demacia/datacenter/position/errors"
	"demacia/datacenter/position/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type PositionEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPositionEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) PositionEditLogic {
	return PositionEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PositionEditLogic) PositionEdit(req types.PositionEditReq) (*types.SuccessReply, error) {

	if req.Id == 0 || req.PositionName == "" {
		return &types.SuccessReply{Success: false}, errlist.InvalidParam
	}

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return &types.SuccessReply{Success: false}, errlist.NoAuth
	}

	positionVal, err := l.svcCtx.PositionModel.SelectById(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.SuccessReply{Success: false}, errors.PositionNotExist
		} else {
			l.Logger.Errorf("select position err:%s", err.Error())
			return &types.SuccessReply{Success: false}, errlist.Unknown
		}
	}
	count, _ := l.svcCtx.PositionModel.SelectExistByPositionName(oid, req.PositionName)
	if count > 0 {
		return &types.SuccessReply{Success: false}, errors.PositionExistErr
	}

	position := &model.Position{PositionName: req.PositionName, Id: positionVal.Id}
	if req.ClassId == 0 {
		position.ClassId = 0
		position.ClassName = ""
	} else {
		positionInfo, _ := l.svcCtx.PositionModel.SelectByClassId(req.ClassId)
		if positionInfo.PositionName != "" {
			return &types.SuccessReply{Success: false}, errors.PositionBindClassErr
		}

		classInfo, err := l.svcCtx.ClassRpc.GetClassInfoById(l.ctx, &class.IdReq{Id: req.ClassId})
		if err != nil {
			return &types.SuccessReply{Success: false}, errlist.Unknown
		}
		position.ClassId = classInfo.Id
		position.ClassName = classInfo.FullName
	}

	if err := l.svcCtx.PositionModel.UpdatePosition(position); err != nil {
		return &types.SuccessReply{Success: false}, errlist.Unknown
	}

	if len(req.DeviceIds) > 0 {
		devicesVal, err := l.svcCtx.PositionDeviceModel.SelectByDeviceIds(req.DeviceIds)
		if err != nil && err != sql.ErrNoRows {
			l.Logger.Errorf("select position_device err:%s", err.Error())
		} else {
			for _, device := range devicesVal {
				if device.PositionId != req.Id {
					return &types.SuccessReply{Success: false}, errors.PositionBindDeviceErr
				}
			}
		}
		devices := make([]*model.PositionDevice, 0, len(req.DeviceIds))
		for _, id := range req.DeviceIds {
			_device, err := l.svcCtx.DeviceRpc.GetDeviceInfoById(l.ctx, &device.IdReq{Id: id})
			if err != nil {
				l.Logger.Errorf("call deviceRpc error:%s", err.Error())
				continue
			}
			devices = append(devices, &model.PositionDevice{
				PositionId: req.Id,
				DeviceId:   id,
				DeviceName: _device.Title,
			})
		}

		if err := l.svcCtx.PositionDeviceModel.DeleteByPid(req.Id); err == nil {
			if err := l.svcCtx.PositionDeviceModel.UpdatePositionDevice(devices, req.Id); err != nil {
				return &types.SuccessReply{Success: false}, errors.PositionEditDeviceErr
			}
		}

	}

	return &types.SuccessReply{Success: true}, nil
}
