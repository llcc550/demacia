// Code generated by goctl. DO NOT EDIT!
// Source: member.proto

package server

import (
	"context"

	"demacia/datacenter/member/rpc/internal/logic"
	"demacia/datacenter/member/rpc/internal/svc"
	"demacia/datacenter/member/rpc/member"
)

type MemberServer struct {
	svcCtx *svc.ServiceContext
	member.UnimplementedMemberServer
}

func NewMemberServer(svcCtx *svc.ServiceContext) *MemberServer {
	return &MemberServer{
		svcCtx: svcCtx,
	}
}

func (s *MemberServer) FindOneById(ctx context.Context, in *member.IdReq) (*member.MemberInfo, error) {
	l := logic.NewFindOneByIdLogic(ctx, s.svcCtx)
	return l.FindOneById(in)
}

func (s *MemberServer) FindOneByUserName(ctx context.Context, in *member.UserNameReq) (*member.MemberInfo, error) {
	l := logic.NewFindOneByUserNameLogic(ctx, s.svcCtx)
	return l.FindOneByUserName(in)
}

func (s *MemberServer) Insert(ctx context.Context, in *member.InsertReq) (*member.IdReq, error) {
	l := logic.NewInsertLogic(ctx, s.svcCtx)
	return l.Insert(in)
}

func (s *MemberServer) DeleteById(ctx context.Context, in *member.IdReq) (*member.NullResp, error) {
	l := logic.NewDeleteByIdLogic(ctx, s.svcCtx)
	return l.DeleteById(in)
}