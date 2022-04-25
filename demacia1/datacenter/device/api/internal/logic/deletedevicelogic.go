package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"

	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type DeleteDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteDeviceLogic {
	return DeleteDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDeviceLogic) DeleteDevice(req types.DeviceIdReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	deviceInfo, err := l.svcCtx.DeviceModel.GetDeviceInfoById(req.DeviceId)
	if err != nil || deviceInfo.OrgId != orgId {
		return errlist.NoAuth
	}
	err = l.svcCtx.DeviceModel.DeleteOneById(req.DeviceId)
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		err = l.svcCtx.DeviceGroupModel.DeleteByDeviceId(req.DeviceId)
		if err != nil {
			logx.Errorf("delete device-group err:%s,err:%v", err.Error(), err)
		}
		s := datacenter.Marshal(datacenter.Device, req.DeviceId, datacenter.Delete, nil)
		_ = l.svcCtx.KqPusher.Push(s)
	})
	return nil
}
