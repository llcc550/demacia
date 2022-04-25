package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"
	"demacia/datacenter/device/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
	"sort"
)

type SelectorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelectorLogic(ctx context.Context, svcCtx *svc.ServiceContext) SelectorLogic {
	return SelectorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectorLogic) Selector() (resp *types.Selector, err error) {
	resp = &types.Selector{List: []*types.DeviceGroup{}}
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, err
	}
	//根据学校id查询所有设备组
	groups, err := l.svcCtx.GroupModel.FindListByOrgId(orgId)
	if err != nil {
		return resp, errlist.DeviceNotExist
	}
	//根据学校id获取所有设备
	groups = append(groups, &model.Group{
		Id:    0,
		OrgId: orgId,
		Name:  "未分组",
	})

	results, err := mr.MapReduce(func(source chan<- interface{}) {
		for _, item := range groups {
			source <- item
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		each := item.(*model.Group)
		groupInfo := &types.DeviceGroup{
			Id:       0,
			Name:     "",
			Children: []*types.Children{},
		}
		groupInfo.Id = each.Id
		groupInfo.Name = each.Name

		//获取组下面所有设备
		deviceList, err := l.svcCtx.DeviceGroupModel.FindListByGroupId([]int64{each.Id})
		if err != nil {
			logx.Errorf("get group device error by id,groupId:%v,error:%+v", each.Id, err)
			cancel(err)
			return
		}

		for _, v := range *deviceList {
			device := &types.Children{
				Id:   0,
				Name: "",
				Pid:  0,
			}
			device.Id = v.DeviceId
			device.Pid = v.GroupId
			groupInfo.Children = append(groupInfo.Children, device)
		}
		writer.Write(groupInfo)
	}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		list := make([]*types.DeviceGroup, 0)
		for p := range pipe {
			list = append(list, p.(*types.DeviceGroup))
		}
		writer.Write(list)
	})
	if err != nil {
		logx.Errorf("get group device error by id,error:%s", err.Error())
		return resp, errlist.DeviceNotExist
	}
	//批量获取设备title 获取device_id
	groupsTree := results.([]*types.DeviceGroup)
	deviceIdList := make([]int64, 0)
	for _, item := range groupsTree {
		for _, child := range item.Children {
			deviceIdList = append(deviceIdList, child.Id)
		}
	}
	//批量获取设备title
	device, err := l.svcCtx.DeviceModel.ListByIds(deviceIdList)
	if err != nil {
		return resp, errlist.DeviceNotExist
	}
	deviceMap := map[int64]*model.Device{}
	for _, v := range device {
		deviceMap[v.Id] = &model.Device{
			Id:      v.Id,
			Sn:      v.Sn,
			OrgId:   v.OrgId,
			Title:   v.Title,
			Network: v.Network,
		}
	}

	for _, group := range groupsTree {
		for _, child := range group.Children {
			deviceTitle := ""
			if _, ok := deviceMap[child.Id]; ok {
				deviceTitle = deviceMap[child.Id].Title
			}
			child.Name = deviceTitle
		}
	}
	//sort 排序
	sortDeviceGroupList := types.DeviceGroupList{}
	sortDeviceGroupList = groupsTree

	for _, v := range sortDeviceGroupList {
		childrenList := types.ChildrenList{}
		childrenList = append(childrenList, v.Children...)
		sort.Sort(childrenList)
		v.Children = childrenList
	}
	sort.Sort(sortDeviceGroupList)
	resp.List = sortDeviceGroupList
	return resp, nil
}
