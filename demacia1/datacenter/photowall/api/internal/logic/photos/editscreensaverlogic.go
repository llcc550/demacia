package photos

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/device/rpc/device"
	"demacia/datacenter/photowall/common"
	"demacia/datacenter/photowall/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"

	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type EditscreensaverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditscreensaverLogic(ctx context.Context, svcCtx *svc.ServiceContext) EditscreensaverLogic {
	return EditscreensaverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Editscreensaver 编辑屏保
func (l *EditscreensaverLogic) Editscreensaver(req types.EditReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return errlist.NoAuth
	}
	if req.PhotoId == 0 && req.DeviceId == 0 || req.PhotoId > 0 && req.DeviceId > 0 {
		return errlist.InvalidParam
	}
	if req.PhotoId > 0 {
		_ = l.EditScreenSaverByPhotoId(orgId, req)
	} else if req.DeviceId > 0 {
		_ = l.EditScreenSaverByDevice(orgId, req)
	}
	return nil
}

// EditScreenSaverByPhotoId 相册编辑屏保
func (l *EditscreensaverLogic) EditScreenSaverByPhotoId(orgId int64, req types.EditReq) error {
	threading.GoSafe(func() {
		// 清除资源所有设备的屏保时间
		err := l.svcCtx.DevicePhotoModel.ClearTimeByPhotoId(common.ScreenSaver, req.PhotoId)
		if err != nil {
			return
		}
		for _, deviceId := range req.DeviceIds {
			// 查询设备是否存在
			deviceInfo, err := l.svcCtx.DeviceRpc.GetDeviceInfoById(l.ctx, &device.IdReq{Id: deviceId})
			if err != nil {
				l.Logger.Errorf("photowall EditScreenSaverByFolder use Device[Rpc] err:%s", err.Error())
				continue
			}
			// 查询设备和资源关联关系
			devicephotoInfo, err := l.svcCtx.DevicePhotoModel.FindPhotoByPhotoIdAndDeviceId(&model.DevicePhoto{
				OrgId:    orgId,
				DeviceId: deviceInfo.Id,
				PhotoId:  req.PhotoId,
			})
			// 不存在关联则新增
			if devicephotoInfo == nil {
				deviceBindPhoto := model.DevicePhoto{
					OrgId:                orgId,
					DeviceId:             req.DeviceId,
					PhotoId:              req.PhotoId,
					ScreensaverStartTime: sql.NullString{String: req.ScreenSaverStartTime},
					ScreensaverEndTime:   sql.NullString{String: req.ScreenSaverEndTime},
					ScreenSaverWaitTime:  req.ScreenSaverWaitTime,
				}
				deviceBindPhoto.ScreensaverStartTime.String = req.ScreenSaverStartTime
				deviceBindPhoto.ScreensaverEndTime.String = req.ScreenSaverEndTime
				_, _ = l.svcCtx.DevicePhotoModel.Insert(&deviceBindPhoto)
			} else {
				// 存在则修改
				_ = l.svcCtx.DevicePhotoModel.UpdateTimeById(&model.DevicePhoto{
					OrgId:                orgId,
					PhotoId:              req.PhotoId,
					ScreensaverStartTime: sql.NullString{String: req.ScreenSaverStartTime},
					ScreensaverEndTime:   sql.NullString{String: req.ScreenSaverEndTime},
					ScreenSaverWaitTime:  req.ScreenSaverWaitTime,
				})
			}
			continue
		}
	})
	return nil
}

// EditScreenSaverByDevice 设备编辑屏保
func (l *EditscreensaverLogic) EditScreenSaverByDevice(orgId int64, req types.EditReq) error {
	return nil
}
