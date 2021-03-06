// Code generated by goctl. DO NOT EDIT!
// Source: databus.proto

package server

import (
	"context"

	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/databus/rpc/internal/logic"
	"demacia/datacenter/databus/rpc/internal/svc"
)

type DatabusServer struct {
	svcCtx *svc.ServiceContext
	databus.UnimplementedDatabusServer
}

func NewDatabusServer(svcCtx *svc.ServiceContext) *DatabusServer {
	return &DatabusServer{
		svcCtx: svcCtx,
	}
}

func (s *DatabusServer) Create(ctx context.Context, in *databus.Req) (*databus.Res, error) {
	l := logic.NewCreateLogic(ctx, s.svcCtx)
	return l.Create(in)
}

func (s *DatabusServer) Update(ctx context.Context, in *databus.Req) (*databus.Res, error) {
	l := logic.NewUpdateLogic(ctx, s.svcCtx)
	return l.Update(in)
}

func (s *DatabusServer) Delete(ctx context.Context, in *databus.Req) (*databus.Res, error) {
	l := logic.NewDeleteLogic(ctx, s.svcCtx)
	return l.Delete(in)
}

func (s *DatabusServer) Log(ctx context.Context, in *databus.LogReq) (*databus.Res, error) {
	l := logic.NewLogLogic(ctx, s.svcCtx)
	return l.Log(in)
}
