package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"

	"demacia/common/datacenter"
	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type UpdateTitleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTitleLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateTitleLogic {
	return UpdateTitleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTitleLogic) UpdateTitle(req types.UpdateTitleReq) error {
	// todo: 参数合法性验证
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	err = mr.Finish(func() error {
		deviceInfo, err := l.svcCtx.DeviceModel.GetDeviceInfoById(req.DeviceId)
		if err != nil || deviceInfo.OrgId != orgId {
			return errlist.NoAuth
		}
		return nil
	}, func() error {
		//设备title唯一
		deviceInfoByName, err := l.svcCtx.DeviceModel.GetOneByTitle(orgId, req.DeviceTitle)
		if err == nil && deviceInfoByName.Id != req.DeviceId {
			return errlist.DeviceTitleExist
		}
		return nil
	})

	err = l.svcCtx.DeviceModel.UpdateTitle(req.DeviceId, req.DeviceTitle)
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		s := datacenter.Marshal(datacenter.Device, req.DeviceId, datacenter.Update, datacenter.OrganizationData{Title: req.DeviceTitle})
		_ = l.svcCtx.KqPusher.Push(s)
	})

	return nil
}
