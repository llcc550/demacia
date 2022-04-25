package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/datacenter/device/model"

	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GroupListLogic {
	return GroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupListLogic) GroupList(req types.GroupListReq) (resp *types.GroupList, err error) {
	resp = &types.GroupList{List: []types.GroupDetail{}}
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, err
	}
	groupList, err := l.svcCtx.GroupModel.FindListByGroupName(orgId, req.GroupName)
	if err != nil {
		return resp, err
	}
	//遍历组的信息 生成groupId切片，通过切片去查询关联关系
	groupIds := make([]int64, 0)
	for _, v := range groupList {
		groupIds = append(groupIds, v.Id)
	}
	deviceGroup, err := l.svcCtx.DeviceGroupModel.FindListByGroupId(groupIds)
	if err != nil {
		for _, v := range groupList {
			resp.List = append(resp.List, types.GroupDetail{
				GroupId:   v.Id,
				GroupName: v.Name,
				Device:    []types.DeviceDetail{},
			})
		}
	}
	deviceIds := make([]int64, 0)
	deviceGroupRelationMap := map[int64][]int64{}
	for _, v := range *deviceGroup {
		deviceIds = append(deviceIds, v.DeviceId)
		deviceGroupRelationMap[v.GroupId] = append(deviceGroupRelationMap[v.GroupId], v.DeviceId)
	}
	deviceIds = append(deviceIds, 0)
	deviceList, _, err := l.svcCtx.DeviceModel.ListByConditions(deviceIds, orgId, "", "", 0, 0, 0)
	if err != nil {
		for _, v := range groupList {
			resp.List = append(resp.List, types.GroupDetail{
				GroupId:   v.Id,
				GroupName: v.Name,
				Device:    []types.DeviceDetail{},
			})
		}
	}
	deviceMap := map[int64]model.Device{}
	for _, v := range deviceList {
		deviceMap[v.Id] = model.Device{
			Id:      v.Id,
			Sn:      v.Sn,
			OrgId:   v.OrgId,
			Title:   v.Title,
			Network: v.Network,
		}

	}
	for _, group := range groupList {
		deviceSlice := make([]types.DeviceDetail, 0)
		for _, deviceId := range deviceGroupRelationMap[group.Id] {
			deviceSlice = append(deviceSlice, types.DeviceDetail{
				DeviceId:    deviceId,
				DeviceTitle: deviceMap[deviceId].Title,
				Sn:          deviceMap[deviceId].Sn,
				Group:       []types.Group{},
				Network:     deviceMap[deviceId].Network,
			})
		}
		resp.List = append(resp.List, types.GroupDetail{
			GroupId:   group.Id,
			GroupName: group.Name,
			Device:    deviceSlice,
		})
	}

	return resp, nil
}
