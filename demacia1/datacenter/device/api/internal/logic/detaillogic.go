package logic

import (
	"context"
	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) DetailLogic {
	return DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req types.DeviceIdReq) (resp *types.DeviceDetail, err error) {
	resp = &types.DeviceDetail{
		DeviceId:    0,
		DeviceTitle: "",
		Sn:          "",
		Group:       []types.Group{},
		Network:     0,
	}
	groupSlice := make([]types.Group, 0)
	deviceInfo, err := l.svcCtx.DeviceModel.GetDeviceInfoById(req.DeviceId)
	if err != nil {
		return resp, err
	}
	deviceGroupList, err := l.svcCtx.DeviceGroupModel.FindListByDeviceId([]int64{req.DeviceId})
	if err != nil {
		return nil, err
	}
	groupIds := make([]int64, 0)
	for _, v := range *deviceGroupList {
		groupIds = append(groupIds, v.GroupId)
	}
	groupList, err := l.svcCtx.GroupModel.FindListByIds(groupIds)
	if err != nil {
		return resp, err
	}
	for _, v := range groupList {
		groupSlice = append(groupSlice, types.Group{
			GroupId:   v.Id,
			GroupName: v.Name,
		})
	}
	resp.Sn = deviceInfo.Sn
	resp.Group = groupSlice
	resp.DeviceId = deviceInfo.Id
	resp.DeviceTitle = deviceInfo.Title
	resp.Network = deviceInfo.Network
	return resp, nil
}
