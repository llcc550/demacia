package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"
	"demacia/datacenter/device/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type InsertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) InsertLogic {
	return InsertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InsertLogic) Insert(req types.InsertReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	err = mr.Finish(func() error {
		_, err = l.svcCtx.DeviceModel.GetOneByTitle(orgId, req.DeviceTitle)
		if err == nil {
			return errlist.DeviceTitleExist
		}
		return nil
	}, func() error {
		_, err = l.svcCtx.DeviceModel.GetOneBySn(req.Sn)
		if err == nil {
			return errlist.DeviceSnExist
		}
		return nil
	})
	if err != nil {
		return err
	}
	deviceId, err := l.svcCtx.DeviceModel.InsertOne(&model.Device{
		Sn:    req.Sn,
		OrgId: req.OrgId,
		Title: req.DeviceTitle,
	})
	if err != nil {
		return errlist.InvalidParam
	}
	threading.GoSafe(func() {
		s := datacenter.Marshal(datacenter.Device, deviceId, datacenter.Add, datacenter.OrganizationData{Title: req.DeviceTitle})
		_ = l.svcCtx.KqPusher.Push(s)
	})
	return nil
}
