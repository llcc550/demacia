package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"
	"demacia/datacenter/device/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListLogic {
	return ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req types.ListReq) (resp *types.ListResponse, err error) {
	resp = &types.ListResponse{
		List:  []types.DeviceDetail{},
		Count: 0,
	}
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, err
	}
	//查询学校下所以组 生成map
	groupInfoMap := map[int64]model.Group{}
	group, err := l.svcCtx.GroupModel.FindListByOrgId(orgId)
	if err != nil {
		return resp, err
	}
	for _, v := range group {
		groupInfoMap[v.Id] = model.Group{
			Id:    v.Id,
			OrgId: v.OrgId,
			Name:  v.Name,
		}
	}

	//根据条件查询出所有设备
	deviceIds := make([]int64, 0)
	if req.GroupId != 0 {
		deviceGroup, err := l.svcCtx.DeviceGroupModel.FindListByGroupId([]int64{req.GroupId})
		if err == nil {
			for _, v := range *deviceGroup {
				deviceIds = append(deviceIds, v.DeviceId)
			}
		}
	}
	deviceIds = append(deviceIds, 0)
	deviceList, count, err := l.svcCtx.DeviceModel.ListByConditions(deviceIds, orgId, req.DeviceTitle, req.Sn, req.Network, req.Page, req.Limit)
	if err != nil || count == 0 {
		return resp, err
	}

	//根据设备id去查询组的关系
	getDeviceIds := make([]int64, 0)
	for _, v := range deviceList {
		getDeviceIds = append(getDeviceIds, v.Id)
	}

	//查询关系 放到map
	deviceGroupRelationMap := map[int64][]int64{}
	groupList, err := l.svcCtx.DeviceGroupModel.FindListByDeviceId(getDeviceIds)
	if err != nil {
		return resp, err
	}
	for _, v := range *groupList {
		deviceGroupRelationMap[v.DeviceId] = append(deviceGroupRelationMap[v.DeviceId], v.GroupId)
	}
	groupSlice := make([]types.Group, 0)
	for _, device := range deviceList {
		for _, v := range deviceGroupRelationMap[device.Id] {
			groupSlice = append(groupSlice, types.Group{
				GroupId:   groupInfoMap[v].Id,
				GroupName: groupInfoMap[v].Name,
			})
		}
		resp.List = append(resp.List, types.DeviceDetail{
			DeviceId:    device.Id,
			DeviceTitle: device.Title,
			Sn:          device.Sn,
			Group:       groupSlice,
			Network:     device.Network,
		})
	}
	resp.Count = count
	return resp, nil
}
