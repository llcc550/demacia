package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"

	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GroupDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) GroupDetailLogic {
	return GroupDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupDetailLogic) GroupDetail(req types.GroupIdReq) (resp *types.GroupDetail, err error) {
	resp = &types.GroupDetail{
		GroupId:   0,
		GroupName: "",
		Device:    []types.DeviceDetail{},
	}
	deviceSlice := make([]types.DeviceDetail, 0)
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, err
	}
	groupInfo, err := l.svcCtx.GroupModel.FindOne(req.GroupId)
	if err != nil || groupInfo.OrgId != orgId {
		return resp, errlist.NoAuth
	}
	//去关联表查询组下面的设备id
	deviceIds := make([]int64, 0)
	deviceGroupList, err := l.svcCtx.DeviceGroupModel.FindListByGroupId([]int64{req.GroupId})
	if err != nil {
		resp.Device = deviceSlice
	}
	for _, v := range *deviceGroupList {
		deviceIds = append(deviceIds, v.DeviceId)
	}
	deviceIds = append(deviceIds, 0)
	//根据device_id查询device信息
	deviceList, _, err := l.svcCtx.DeviceModel.ListByConditions(deviceIds, orgId, "", "", 0, 0, 0)
	if err != nil {
		resp.Device = deviceSlice
	}
	for _, v := range deviceList {
		deviceSlice = append(deviceSlice, types.DeviceDetail{
			DeviceId:    v.Id,
			DeviceTitle: v.Title,
			Sn:          v.Sn,
			Group:       []types.Group{},
			Network:     v.Network,
		})
	}
	resp.Device = deviceSlice
	resp.GroupId = groupInfo.Id
	resp.GroupName = groupInfo.Name
	return resp, nil
}
