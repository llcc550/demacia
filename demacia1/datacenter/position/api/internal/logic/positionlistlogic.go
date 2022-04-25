package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"

	"demacia/datacenter/position/api/internal/svc"
	"demacia/datacenter/position/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type PositionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPositionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) PositionListLogic {
	return PositionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PositionListLogic) PositionList(req types.PageReq) (*types.PositionListReply, error) {
	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return &types.PositionListReply{Positions: []*types.Position{}}, errlist.NoAuth
	}
	if req.Page < 0 {
		req.Page = 0
	}
	if req.Limit < 0 {
		req.Limit = 10
	}
	positionList, count, err := l.svcCtx.PositionModel.SelectList(oid, req.Page, req.Limit, req.PositionName)

	if err != nil {
		return &types.PositionListReply{Positions: []*types.Position{}}, errlist.Unknown
	}

	if len(positionList) == 0 {
		return &types.PositionListReply{Positions: []*types.Position{}}, nil
	}

	resp := types.PositionListReply{}
	resp.Positions = []*types.Position{}
	resp.Count = count
	ids := make([]int64, 0, len(positionList))

	for _, position := range positionList {
		ids = append(ids, position.Id)
		resp.Positions = append(resp.Positions, &types.Position{
			Id:           position.Id,
			PositionName: position.PositionName,
			ClassName:    position.ClassName,
			Devices:      []*types.Device{},
		})
	}

	deviceList, err := l.svcCtx.PositionDeviceModel.SelectByPidList(ids)
	for _, device := range deviceList {
		for _, position := range resp.Positions {
			if device.PositionId == position.Id {
				position.Devices = append(position.Devices, &types.Device{
					Id:         device.DeviceId,
					DeviceName: device.DeviceName,
				})
			}
		}
	}

	return &resp, nil
}
