// Code generated by goctl. DO NOT EDIT!
// Source: organization.proto

package server

import (
	"context"

	"demacia/datacenter/organization/rpc/internal/logic"
	"demacia/datacenter/organization/rpc/internal/svc"
	"demacia/datacenter/organization/rpc/organization"
)

type OrganizationServer struct {
	svcCtx *svc.ServiceContext
	organization.UnimplementedOrganizationServer
}

func NewOrganizationServer(svcCtx *svc.ServiceContext) *OrganizationServer {
	return &OrganizationServer{
		svcCtx: svcCtx,
	}
}

func (s *OrganizationServer) FindOne(ctx context.Context, in *organization.IdReply) (*organization.OrgInfo, error) {
	l := logic.NewFindOneLogic(ctx, s.svcCtx)
	return l.FindOne(in)
}
