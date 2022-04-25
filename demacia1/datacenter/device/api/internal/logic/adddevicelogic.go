package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/basefunc"
	"demacia/common/errlist"
	"demacia/datacenter/device/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"

	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type AddDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddDeviceLogic {
	return AddDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddDeviceLogic) AddDevice(req types.AddDeviceReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	groupInfo, err := l.svcCtx.GroupModel.FindOne(req.GroupId)
	if err != nil || groupInfo.OrgId != orgId {
		return errlist.NoAuth
	}
	deviceIds := basefunc.RemoveRepByLoop(req.DeviceIds)
	err = l.svcCtx.DeviceGroupModel.DeleteByGroupId(req.GroupId)
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		for _, id := range deviceIds {
			_, err := l.svcCtx.DeviceGroupModel.Insert(&model.DeviceGroup{
				Id:       0,
				DeviceId: id,
				GroupId:  req.GroupId,
			})
			if err != nil {
				logx.Errorf("insert device_group err:%s", err.Error())
				continue
			}
		}
	},
	)

	return nil
}
