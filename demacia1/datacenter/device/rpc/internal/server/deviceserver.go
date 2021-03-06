// Code generated by goctl. DO NOT EDIT!
// Source: device.proto

package server

import (
	"context"

	"demacia/datacenter/device/rpc/device"
	"demacia/datacenter/device/rpc/internal/logic"
	"demacia/datacenter/device/rpc/internal/svc"
)

type DeviceServer struct {
	svcCtx *svc.ServiceContext
}

func NewDeviceServer(svcCtx *svc.ServiceContext) *DeviceServer {
	return &DeviceServer{
		svcCtx: svcCtx,
	}
}

func (s *DeviceServer) GetDeviceInfoById(ctx context.Context, in *device.IdReq) (*device.DeviceInfo, error) {
	l := logic.NewGetDeviceInfoByIdLogic(ctx, s.svcCtx)
	return l.GetDeviceInfoById(in)
}
