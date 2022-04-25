// Code generated by goctl. DO NOT EDIT!
// Source: common.proto

package commonclient

import (
	"context"

	"demacia/datacenter/common/rpc/common"

	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AreaReq     = common.AreaReq
	AreaResp    = common.AreaResp
	HolidayReq  = common.HolidayReq
	HolidayResp = common.HolidayResp

	Common interface {
		FindAreaInfo(ctx context.Context, in *AreaReq, opts ...grpc.CallOption) (*AreaResp, error)
		FindHolidayInfo(ctx context.Context, in *HolidayReq, opts ...grpc.CallOption) (*HolidayResp, error)
	}

	defaultCommon struct {
		cli zrpc.Client
	}
)

func NewCommon(cli zrpc.Client) Common {
	return &defaultCommon{
		cli: cli,
	}
}

func (m *defaultCommon) FindAreaInfo(ctx context.Context, in *AreaReq, opts ...grpc.CallOption) (*AreaResp, error) {
	client := common.NewCommonClient(m.cli.Conn())
	return client.FindAreaInfo(ctx, in, opts...)
}

func (m *defaultCommon) FindHolidayInfo(ctx context.Context, in *HolidayReq, opts ...grpc.CallOption) (*HolidayResp, error) {
	client := common.NewCommonClient(m.cli.Conn())
	return client.FindHolidayInfo(ctx, in, opts...)
}