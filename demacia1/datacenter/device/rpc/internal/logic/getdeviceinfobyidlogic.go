package logic

import (
	"context"

	"demacia/datacenter/device/errors"
	"demacia/datacenter/device/rpc/device"
	"demacia/datacenter/device/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetDeviceInfoByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceInfoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceInfoByIdLogic {
	return &GetDeviceInfoByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeviceInfoByIdLogic) GetDeviceInfoById(in *device.IdReq) (*device.DeviceInfo, error) {
	info, err := l.svcCtx.DeviceModel.GetDeviceInfoById(in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, errors.DeviceNotExist.Error())
	}
	return &device.DeviceInfo{Id: info.Id, Title: info.Title}, nil
}
